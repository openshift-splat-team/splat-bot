package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/prstatus"

	splathub "github.com/openshift-splat-team/splat-bot/pkg/github"

	"github.com/beatlabs/github-auth/app"
	"github.com/beatlabs/github-auth/key"
	githubql "github.com/shurcooL/githubv4"
)

const (
	githubAppId = "858938"
)

var (
	prList   []prstatus.PullRequest
	githubID string

	boldStyle = slack.RichTextSectionTextStyle{Bold: true}
)

var PullRequestAttributes = data.Attributes{
	Commands:       []string{"pull-requests"},
	RequireMention: true,
	Callback: func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		var err error
		prList, err = fetchPullRequests(args)

		if err != nil {
			return nil, fmt.Errorf("user not allowed: %v", err)
		}

		return generateOutput()
	},
	AllowNonSplatUsers:  true,
	RequiredArgs:        2,
	HelpMarkdown:        "retrieve list of pull requests open for the specified user: `pull-requests [user]`",
	ResponseIsEphemeral: true,
	RespondInChannel:    true,
	ShouldMatch: []string{
		"pull-requests rvanderp3",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}

func getGithubKeyPath() string {
	keyPath := os.Getenv("GITHUB_KEY_PATH")
	if keyPath == "" {
		keyPath = "data/private.key"
	}
	return keyPath
}

func init() {
	keyPath := getGithubKeyPath()
	_, err := os.ReadFile(keyPath)
	if err != nil {
		log.Printf("error reading file %s: %v", keyPath, err)
	}
	// If key not found, disable command and log missing file.
	if err != nil {
		fmt.Printf("error loading knowledge entries: %v", err)
		fmt.Println("Skipping adding of knowledge-based actions.")
		return
	}
	AddCommand(PullRequestAttributes)
}

func generateOutput() ([]slack.MsgOption, error) {

	var messageBlocks []slack.Block
	truncated := false
	log.Printf("Attempting to creating %v PR entries.", len(prList))
	//var prResultsBuffer strings.Builder
	for index, pr := range prList {
		if len(messageBlocks)+3+3 > 50 {
			log.Printf("Due to number of blocks, stopping at index: %d", index)
			truncated = true
			break
		}
		// Generate divider after first PR
		if index > 0 {
			divider := slack.NewDividerBlock()
			messageBlocks = append(messageBlocks, divider)
		}

		// Master block for as much text as possible.
		prBlock := slack.NewRichTextBlock("")

		// Generate Header (title)
		createFieldWithValue(prBlock, "Title: ", string(pr.Title), true)

		// Generate Project
		createFieldWithValue(prBlock, "Project: ", string(pr.Repository.Name), false)

		// Generate Labels
		prLabelsText := slack.NewRichTextSectionTextElement("Labels: ", &boldStyle)
		prLabelsSection := slack.NewRichTextSection(prLabelsText)
		prBlock.Elements = append(prBlock.Elements, prLabelsSection)

		if len(pr.Labels.Nodes) > 0 {
			var prLabelBullets []slack.RichTextElement
			for _, label := range pr.Labels.Nodes {
				prLabelListEntry := slack.NewRichTextSectionTextElement(string(label.Label.Name), nil)
				prLabelListSection := slack.NewRichTextSection(prLabelListEntry)
				prLabelBullets = append(prLabelBullets, prLabelListSection)
			}
			prLabelsList := slack.NewRichTextList("bullet", 0, prLabelBullets...)
			prBlock.Elements = append(prBlock.Elements, prLabelsList)
		} else {
			prLabelsNone := slack.NewRichTextSectionTextElement("None", nil)
			prLabelsSection.Elements = append(prLabelsSection.Elements, prLabelsNone)
		}

		// Generate Mergeable info
		createFieldWithValue(prBlock, "Merge State: ", string(pr.Mergeable), false)

		// Add the PR Block
		messageBlocks = append(messageBlocks, prBlock)

		// Generate button to open PR
		openPrText := slack.NewTextBlockObject("plain_text", "View PR", true, false)
		openPr := slack.NewButtonBlockElement("", "", openPrText)
		openPr.URL = generatePrURL(pr) //"https://www.google.com"
		openPr.Style = "primary"
		openPrActionBlock := slack.NewActionBlock("", openPr)
		messageBlocks = append(messageBlocks, openPrActionBlock)
	}

	// Add block w/ button to close ephemeral message
	divider := slack.NewDividerBlock()
	messageBlocks = append(messageBlocks, divider)

	if truncated {
		// TODO: Need to add output to notify user of truncation
		log.Print("Detected truncated results.")
	}

	lineReturn := slack.NewRichTextSectionTextElement("\n", nil)
	lineReturnSection := slack.NewRichTextSection(lineReturn)
	messageBlocks = append(messageBlocks, slack.NewRichTextBlock("", lineReturnSection))

	closeText := slack.NewTextBlockObject("plain_text", "Close", true, false)
	closeButton := slack.NewButtonBlockElement("", "", closeText)
	closeActionBlock := slack.NewActionBlock("", closeButton)
	messageBlocks = append(messageBlocks, closeActionBlock)

	log.Printf("Number of blocks: %d", len(messageBlocks))

	buffer := bytes.NewBuffer([]byte{})
	msg := slack.Msg{
		Blocks: slack.Blocks{BlockSet: messageBlocks},
	}
	if err := json.NewEncoder(buffer).Encode(msg); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Print(buffer.String())
	}

	return []slack.MsgOption{
		slack.MsgOptionBlocks(messageBlocks...),
	}, nil
}

func createFieldWithValue(block *slack.RichTextBlock, fieldText, fieldValue string, addLineReturn bool) {
	fieldLabel := slack.NewRichTextSectionTextElement(fieldText, &boldStyle)
	fieldTextElement := slack.NewRichTextSectionTextElement(fieldValue, nil)
	fieldSection := slack.NewRichTextSection(fieldLabel, fieldTextElement)
	block.Elements = append(block.Elements, fieldSection)

	if addLineReturn {
		lineReturn := slack.NewRichTextSectionTextElement("\n", nil)
		lineReturnSection := slack.NewRichTextSection(lineReturn)
		block.Elements = append(block.Elements, lineReturnSection)
	}
}

type GithubTokenResponse struct {
	Token   string `json:"token,omitempty"`
	Expires string `json:"expires_at,omitempty"`
	Repos   string `json:"repository_selection,omitempty"`
}

func generatePrURL(pr prstatus.PullRequest) string {
	return fmt.Sprintf("https://github.com/%v/%v/pull/%d", pr.Repository.Owner.Login, pr.Repository.Name, pr.Number)
}

func getGithubToken() (string, error) {
	githubID = "858938"

	// load from a file
	keyFile, err := key.FromFile(getGithubKeyPath())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// load from data
	//key, err := key.Parse(bytes)

	// Create an App Config using the App ID and the private key
	githubApp, err := app.NewConfig(githubID, keyFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Get an *http.Client
	client := githubApp.Client()

	// The client can be used to send authenticated requests
	resp, err := client.Get("https://api.github.com/app")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Response: %v\n", resp)

	install, err := githubApp.InstallationConfig("48639702")
	if err != nil {
		fmt.Printf("Unable to get install app id: %v\n", err)
	} else {
		fmt.Printf("Got install app %v\n", install)
	}

	// Generate JWT
	fmt.Printf("Time: %v\n", time.Now().String())
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iss": githubID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(),
			"alg": "RS256",
		})
	tokenString, err := token.SignedString(keyFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Get device code
	//request, err := http.NewRequest("POST", "https://github.com/login/device/code?client_id=Iv1.3d0af71eca6ada9f", nil)
	//install.Client(context.TODO()).Post("https://api.github.com/app/installations/48639702/access_tokens")

	// Get Access Token
	var request *http.Request
	fmt.Printf("%v %v\n", "Bearer", tokenString)
	request, err = http.NewRequest("POST", "https://api.github.com/app/installations/48639702/access_tokens", bytes.NewBuffer(nil))
	if err != nil {
		return "", err
	}
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("%v %v", "Bearer", tokenString))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	httpClient := &http.Client{}
	resp, err = httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenResp := GithubTokenResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&tokenResp); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	gitToken := tokenResp.Token

	return gitToken, nil
}

func ConstructSearchQuery(login string) string {
	tokens := []string{"is:pr", "state:open", "author:" + login}
	return strings.Join(tokens, " ")
}

// QueryPullRequests is a query function that returns a list of open pull requests owned by the user whose access token
// is consumed by the github client.
func QueryPullRequests(ctx context.Context, ghc GithubQuerier, query string) ([]prstatus.PullRequest, error) {
	var prs []prstatus.PullRequest
	vars := map[string]interface{}{
		"query":        (githubql.String)(query),
		"searchCursor": (*githubql.String)(nil),
	}
	var totalCost int
	var remaining int
	for {
		sq := searchQuery{}
		if err := ghc.QueryWithGitHubAppsSupport(ctx, &sq, vars, ""); err != nil {
			return nil, err
		}
		totalCost += int(sq.RateLimit.Cost)
		remaining = int(sq.RateLimit.Remaining)
		for _, n := range sq.Search.Nodes {
			org := string(n.PullRequest.Repository.Owner.Login)
			repo := string(n.PullRequest.Repository.Name)
			ref := string(n.PullRequest.HeadRefOID)
			if org == "" || repo == "" || ref == "" {
				// TODO: da.log.Warningf("Skipped empty pull request returned by query \"%s\": %v", query, n.PullRequest)
				continue
			}
			prs = append(prs, n.PullRequest)
		}
		if !sq.Search.PageInfo.HasNextPage {
			break
		}
		vars["searchCursor"] = githubql.NewString(sq.Search.PageInfo.EndCursor)
	}
	// TODO: da.log.Infof("Search for query \"%s\" cost %d point(s). %d remaining.", query, totalCost, remaining)
	fmt.Printf("Search for query \"%s\" cost %d point(s). %d remaining.\n", query, totalCost, remaining)
	return prs, nil
}

func fetchPullRequests(args []string) ([]prstatus.PullRequest, error) {
	var prList []prstatus.PullRequest

	gitToken, err := getGithubToken()
	if err != nil {
		return nil, err
	}

	// Create GithubOptions
	githubOptions := splathub.GitHubOptions{
		Host:              "github.com",
		Endpoint:          splathub.NewStrings(github.DefaultAPIEndpoint),
		GraphqlEndpoint:   github.DefaultGraphQLEndpoint,
		AppID:             githubAppId,
		AppPrivateKeyPath: "data/private.key",
	}

	// Create github client
	clientCreator := func(accessToken string) (prstatus.GitHubClient, error) {
		return githubOptions.GitHubClientWithAccessToken(accessToken)
	}
	githubClient, err := clientCreator(gitToken)
	if err != nil {
		fmt.Printf("Error creating github client: %v\n", err)
		return nil, err
	}
	query := ConstructSearchQuery(args[1])
	prList, err = QueryPullRequests(context.TODO(), githubClient, query)
	if err != nil {
		fmt.Printf("Failed to get PRs: %v\n", err)
		return nil, err
	}
	return prList, nil
}

type GithubQuerier interface {
	QueryWithGitHubAppsSupport(ctx context.Context, q interface{}, vars map[string]interface{}, org string) error
}

type searchQuery struct {
	RateLimit struct {
		Cost      githubql.Int
		Remaining githubql.Int
	}
	Search struct {
		PageInfo struct {
			HasNextPage githubql.Boolean
			EndCursor   githubql.String
		}
		Nodes []struct {
			PullRequest prstatus.PullRequest `graphql:"... on PullRequest"`
		}
	} `graphql:"search(type: ISSUE, first: 100, after: $searchCursor, query: $query)"`
}
