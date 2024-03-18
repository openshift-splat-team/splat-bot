package github

import (
	"context"
	"fmt"
	githubql "github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"k8s.io/test-infra/ghproxy/ghcache"
	"k8s.io/test-infra/prow/config/secret"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/throttle"
	"net/http"
	"sync"
	"time"
)

const (
	acceptNone       = ""
	githubApiVersion = "2022-11-28"

	// MaxRequestTime aborts requests that don't return in 5 mins. Longest graphql
	// calls can take up to 2 minutes. This limit should ensure all successful calls
	// return but will prevent an indefinite stall if GitHub never responds.
	MaxRequestTime = 5 * time.Minute

	DefaultMaxRetries    = 8
	DefaultMax404Retries = 2
	DefaultMaxSleepTime  = 2 * time.Minute
	DefaultInitialDelay  = 2 * time.Second
)

// Strings represents the value of a flag that accept multiple strings.
type Strings struct {
	vals    []string
	beenSet bool
}

// NewStrings returns a Strings struct that defaults to the value of def if left unset.
func NewStrings(def ...string) Strings {
	return Strings{
		vals:    def,
		beenSet: false,
	}
}

// Strings returns the slice of strings set for this value instance.
func (s *Strings) Strings() []string {
	return s.vals
}

type ThrottlerSettings struct {
	HourlyTokens int
	Burst        int
}

// TokenGenerator knows how to generate a token for use in git client calls
type TokenGenerator func(org string) (string, error)

// UserGenerator knows how to identify this user for use in git client calls
type UserGenerator func() (string, error)

// ClientOptions holds options for creating a new client
// TODO: This can probably be removed and borrow from prow.
/*type ClientOptions struct {
	// censor knows how to censor output
	Censor func([]byte) []byte

	// the following fields handle auth
	GetToken      func() []byte
	AppID         string
	AppPrivateKey func() *rsa.PrivateKey

	// the following fields determine which server we talk to
	GraphqlEndpoint string
	Bases           []string

	// the following fields determine client retry behavior
	MaxRequestTime, InitialDelay, MaxSleepTime time.Duration
	MaxRetries, Max404Retries                  int

	DryRun bool
	// BaseRoundTripper is the last RoundTripper to be called. Used for testing, gets defaulted to http.DefaultTransport
	BaseRoundTripper http.RoundTripper
}

func (o ClientOptions) Default() ClientOptions {
	if o.MaxRequestTime == 0 {
		o.MaxRequestTime = MaxRequestTime
	}
	if o.InitialDelay == 0 {
		o.InitialDelay = DefaultInitialDelay
	}
	if o.MaxSleepTime == 0 {
		o.MaxSleepTime = DefaultMaxSleepTime
	}
	if o.MaxRetries == 0 {
		o.MaxRetries = DefaultMaxRetries
	}
	if o.Max404Retries == 0 {
		o.Max404Retries = DefaultMax404Retries
	}
	return o
}*/

// GitHubOptions holds options for interacting with GitHub.
//
// Set AllowAnonymous to be true if you want to allow anonymous github access.
// Set AllowDirectAccess to be true if you want to suppress warnings on direct github access (without ghproxy).
type GitHubOptions struct {
	Host              string
	Endpoint          Strings
	GraphqlEndpoint   string
	TokenPath         string
	AllowAnonymous    bool
	AllowDirectAccess bool
	AppID             string
	AppPrivateKeyPath string

	ThrottleHourlyTokens int
	ThrottleAllowBurst   int

	OrgThrottlers       Strings
	ParsedOrgThrottlers map[string]ThrottlerSettings

	// These will only be set after a github client was retrieved for the first time
	TokenGenerator github.TokenGenerator
	UserGenerator  github.UserGenerator

	// the following options determine how the client behaves around retries
	MaxRequestTime time.Duration
	MaxRetries     int
	Max404Retries  int
	InitialDelay   time.Duration
	MaxSleepTime   time.Duration
}

// GitHubClientWithAccessToken creates a GitHub client from an access token.
func (o *GitHubOptions) GitHubClientWithAccessToken(token string) (github.Client, error) {
	options := o.baseClientOptions()
	options.GetToken = func() []byte { return []byte(token) }
	options.AppID = "" // Since we are using a token, we should not use the app auth
	_, _, client, err := NewClientFromOptions(logrus.Fields{}, options)
	return client, err
}

// baseClientOptions populates client options that are derived from flags without processing
func (o *GitHubOptions) baseClientOptions() github.ClientOptions {
	return github.ClientOptions{
		Censor:          secret.Censor,
		AppID:           o.AppID,
		GraphqlEndpoint: o.GraphqlEndpoint,
		Bases:           o.Endpoint.Strings(),
		MaxRequestTime:  o.MaxRequestTime,
		InitialDelay:    o.InitialDelay,
		MaxSleepTime:    o.MaxSleepTime,
		MaxRetries:      o.MaxRetries,
		Max404Retries:   o.Max404Retries,
	}
}

type addHeaderTransport struct {
	upstream http.RoundTripper
}

var userAgentContextKey = &struct{}{}

func (s *addHeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// We have to add this header to enable the Checks scheme preview:
	// https://docs.github.com/en/enterprise-server@2.22/graphql/overview/schema-previews
	// Any GHE version after 2.22 will enable the Checks types per default
	r.Header.Add("Accept", "application/vnd.github.antiope-preview+json")

	// We have to add this header to enable the Merge info scheme preview:
	// https://docs.github.com/en/graphql/overview/schema-previews#merge-info-preview
	r.Header.Add("Accept", "application/vnd.github.merge-info-preview+json")

	// We use the context to pass the UserAgent through the V4 client we depend on
	if v := r.Context().Value(userAgentContextKey); v != nil {
		r.Header.Add("User-Agent", v.(string))
	}

	return s.upstream.RoundTrip(r)
}

// delegate actually does the work to talk to GitHub
type delegate struct {
	time timeClient

	maxRetries    int
	max404Retries int
	maxSleepTime  time.Duration
	initialDelay  time.Duration

	client       httpClient
	bases        []string
	dry          bool
	fake         bool
	usesAppsAuth bool
	throttle     ghThrottler
	getToken     func() []byte
	censor       func([]byte) []byte

	mut      sync.Mutex // protects botName and email
	userData *github.UserData
}

// Interface for how prow interacts with the http client, which we may throttle.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type reloadingTokenSource struct {
	getToken func() []byte
}

// Token is an implementation for oauth2.TokenSource interface.
func (s *reloadingTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: string(s.getToken()),
	}, nil
}

type timeClient interface {
	Sleep(time.Duration)
	Until(time.Time) time.Duration
}

type standardTime struct{}

func (s *standardTime) Sleep(d time.Duration) {
	time.Sleep(d)
}
func (s *standardTime) Until(t time.Time) time.Duration {
	return time.Until(t)
}

// ghThrottler sets a ceiling on the rate of GitHub requests.
// Configure with Client.Throttle().
// It gets reconstructed whenever forUserAgent() is called,
// whereas its *throttle.Throttler remains.
type ghThrottler struct {
	graph gqlClient
	http  httpClient
	*throttle.Throttler
}

func (t *ghThrottler) Do(req *http.Request) (*http.Response, error) {
	org := extractOrgFromContext(req.Context())
	if err := t.Wait(req.Context(), org); err != nil {
		return nil, err
	}
	resp, err := t.http.Do(req)
	if err == nil {
		cacheMode := ghcache.CacheResponseMode(resp.Header.Get(ghcache.CacheModeHeader))
		if ghcache.CacheModeIsFree(cacheMode) {
			// This request was fulfilled by ghcache without using an API token.
			// Refund the throttling token we preemptively consumed.
			logrus.WithFields(logrus.Fields{
				"client":     "github",
				"throttled":  true,
				"cache-mode": string(cacheMode),
			}).Debug("Throttler refunding token for free response from ghcache.")
			t.Refund(org)
		} else {
			logrus.WithFields(logrus.Fields{
				"client":     "github",
				"throttled":  true,
				"cache-mode": string(cacheMode),
				"path":       req.URL.Path,
				"method":     req.Method,
			}).Debug("Used token for request")

		}
	}
	return resp, err
}

func (t *ghThrottler) QueryWithGitHubAppsSupport(ctx context.Context, q interface{}, vars map[string]interface{}, org string) error {
	if err := t.Wait(ctx, extractOrgFromContext(ctx)); err != nil {
		return err
	}
	return t.graph.QueryWithGitHubAppsSupport(ctx, q, vars, org)
}

func (t *ghThrottler) MutateWithGitHubAppsSupport(ctx context.Context, m interface{}, input githubql.Input, vars map[string]interface{}, org string) error {
	if err := t.Wait(ctx, extractOrgFromContext(ctx)); err != nil {
		return err
	}
	return t.graph.MutateWithGitHubAppsSupport(ctx, m, input, vars, org)
}

func (t *ghThrottler) forUserAgent(userAgent string) gqlClient {
	return &ghThrottler{
		graph:     t.graph.forUserAgent(userAgent),
		Throttler: t.Throttler,
	}
}

func extractOrgFromContext(ctx context.Context) string {
	var org string
	if v := ctx.Value(githubOrgHeaderKey); v != nil {
		org = v.(string)
	}
	return org
}

// newReloadingTokenSource creates a reloadingTokenSource.
func newReloadingTokenSource(getToken func() []byte) *reloadingTokenSource {
	return &reloadingTokenSource{
		getToken: getToken,
	}
}

// NewClientFromOptions creates a new client from the options we expose. This method should be used over the more-specific ones.
func NewClientFromOptions(fields logrus.Fields, options github.ClientOptions) (TokenGenerator, UserGenerator, github.Client, error) {
	options = options.Default()

	// Will be nil if github app authentication is used
	if options.GetToken == nil {
		options.GetToken = func() []byte { return nil }
	}
	if options.BaseRoundTripper == nil {
		options.BaseRoundTripper = http.DefaultTransport
	}

	httpClient := &http.Client{
		Transport: options.BaseRoundTripper,
		Timeout:   options.MaxRequestTime,
	}
	graphQLTransport := newAddHeaderTransport(options.BaseRoundTripper)
	c := &client{
		logger: logrus.WithFields(fields).WithField("client", "github"),
		gqlc: &graphQLGitHubAppsAuthClientWrapper{Client: githubql.NewEnterpriseClient(
			options.GraphqlEndpoint,
			&http.Client{
				Timeout: options.MaxRequestTime,
				Transport: &oauth2.Transport{
					Source: newReloadingTokenSource(options.GetToken),
					Base:   graphQLTransport,
				},
			})},
		delegate: &delegate{
			time:          &standardTime{},
			client:        httpClient,
			bases:         options.Bases,
			throttle:      ghThrottler{Throttler: &throttle.Throttler{}},
			getToken:      options.GetToken,
			censor:        options.Censor,
			dry:           options.DryRun,
			usesAppsAuth:  options.AppID != "",
			maxRetries:    options.MaxRetries,
			max404Retries: options.Max404Retries,
			initialDelay:  options.InitialDelay,
			maxSleepTime:  options.MaxSleepTime,
		},
	}
	c.gqlc = c.gqlc.forUserAgent(c.userAgent())

	// Wrap clients with the throttler
	c.wrapThrottler()

	var tokenGenerator func(_ string) (string, error)
	var userGenerator func() (string, error)
	if options.AppID != "" {
		appsTransport, err := newAppsRoundTripper(options.AppID, options.AppPrivateKey, options.BaseRoundTripper, c, options.Bases)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to construct apps auth roundtripper: %w", err)
		}
		httpClient.Transport = appsTransport
		graphQLTransport.upstream = appsTransport

		// Use github apps auth for git actions
		// https://docs.github.com/en/free-pro-team@latest/developers/apps/authenticating-with-github-apps#http-based-git-access-by-an-installation=
		tokenGenerator = func(org string) (string, error) {
			res, _, err := appsTransport.installationTokenFor(org)
			return res, err
		}
		userGenerator = func() (string, error) {
			return "x-access-token", nil
		}
	} else {
		// Use Personal Access token auth for git actions
		tokenGenerator = func(_ string) (string, error) {
			return string(options.GetToken()), nil
		}
		userGenerator = func() (string, error) {
			user, err := c.BotUser()
			if err != nil {
				return "", err
			}
			return user.Login, nil
		}
	}

	return tokenGenerator, userGenerator, c, nil
}

func newAddHeaderTransport(upstream http.RoundTripper) *addHeaderTransport {
	return &addHeaderTransport{upstream}
}
