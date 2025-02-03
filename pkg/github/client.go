package github

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	githubql "github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/version"
)

type key string

const (
	githubOrgHeaderKey key = "X-PROW-GITHUB-ORG"
)

// client interacts with the github api. It is reconstructed whenever
// ForPlugin/ForSubcomment is called to change the Logger and User-Agent
// header, whereas delegate will stay the same.
type client struct {
	// If logger is non-nil, log all method calls with it.
	logger *logrus.Entry
	// identifier is used to add more identification to the user-agent header
	identifier string
	gqlc       gqlClient
	used       bool
	mutUsed    sync.Mutex // protects used
	*delegate
}

func (c *client) EditPullRequest(org, repo string, number int, pr *github.PullRequest) (*github.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetPullRequestDiff(org, repo string, number int) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetPullRequestPatch(org, repo string, number int) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdatePullRequest(org, repo string, number int, title, body *string, open *bool, branch *string, canModify *bool) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetPullRequestChanges(org, repo string, number int) ([]github.PullRequestChange, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListPullRequestComments(org, repo string, number int) ([]github.ReviewComment, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreatePullRequestReviewComment(org, repo string, number int, rc github.ReviewComment) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListReviews(org, repo string, number int) ([]github.Review, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ClosePullRequest(org, repo string, number int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ReopenPullRequest(org, repo string, number int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateReview(org, repo string, number int, r github.DraftReview) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) RequestReview(org, repo string, number int, logins []string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) UnrequestReview(org, repo string, number int, logins []string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) Merge(org, repo string, pr int, details github.MergeDetails) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) IsMergeable(org, repo string, number int, SHA string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListPullRequestCommits(org, repo string, number int) ([]github.RepositoryCommit, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdatePullRequestBranch(org, repo string, number int, expectedHeadSha *string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetRepo(owner, name string) (github.FullRepo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetRepos(org string, isUser bool) ([]github.Repo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetBranches(org, repo string, onlyProtected bool) ([]github.Branch, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetBranchProtection(org, repo, branch string) (*github.BranchProtection, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveBranchProtection(org, repo, branch string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateBranchProtection(org, repo, branch string, config github.BranchProtectionRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddRepoLabel(org, repo, label, description, color string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateRepoLabel(org, repo, label, newName, description, color string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteRepoLabel(org, repo, label string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetRepoLabels(org, repo string) ([]github.Label, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddLabel(org, repo string, number int, label string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddLabelWithContext(ctx context.Context, org, repo string, number int, label string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddLabels(org, repo string, number int, labels ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) AddLabelsWithContext(ctx context.Context, org, repo string, number int, labels ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveLabel(org, repo string, number int, label string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveLabelWithContext(ctx context.Context, org, repo string, number int, label string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) WasLabelAddedByHuman(org, repo string, number int, label string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetFile(org, repo, filepath, commit string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetDirectory(org, repo, dirpath, commit string) ([]github.DirectoryContent, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) IsCollaborator(org, repo, user string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListCollaborators(org, repo string) ([]github.User, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateFork(owner, repo string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) EnsureFork(forkingUser, org, repo string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListRepoTeams(org, repo string) ([]github.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateRepo(owner string, isUser bool, repo github.RepoCreateRequest) (*github.FullRepo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateRepo(owner, name string, repo github.RepoUpdateRequest) (*github.FullRepo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateStatus(org, repo, SHA string, s github.Status) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateStatusWithContext(ctx context.Context, org, repo, SHA string, s github.Status) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListStatuses(org, repo, ref string) ([]github.Status, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetSingleCommit(org, repo, SHA string) (github.RepositoryCommit, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetCombinedStatus(org, repo, ref string) (*github.CombinedStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListCheckRuns(org, repo, ref string) (*github.CheckRunList, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetRef(org, repo, ref string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteRef(org, repo, ref string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListFileCommits(org, repo, path string) ([]github.RepositoryCommit, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateCheckRun(org, repo string, checkRun github.CheckRun) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateIssue(org, repo, title, body string, milestone int, labels, assignees []string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateIssueReaction(org, repo string, id int, reaction string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListIssueComments(org, repo string, number int) ([]github.IssueComment, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListIssueCommentsWithContext(ctx context.Context, org, repo string, number int) ([]github.IssueComment, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetIssueLabels(org, repo string, number int) ([]github.Label, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListIssueEvents(org, repo string, num int) ([]github.ListedIssueEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) AssignIssue(org, repo string, number int, logins []string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) UnassignIssue(org, repo string, number int, logins []string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CloseIssue(org, repo string, number int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CloseIssueAsNotPlanned(org, repo string, number int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ReopenIssue(org, repo string, number int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) FindIssues(query, sort string, asc bool) ([]github.Issue, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) FindIssuesWithOrg(org, query, sort string, asc bool) ([]github.Issue, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListOpenIssues(org, repo string) ([]github.Issue, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetIssue(org, repo string, number int) (*github.Issue, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditIssue(org, repo string, number int, issue *github.Issue) (*github.Issue, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateComment(org, repo string, number int, comment string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateCommentWithContext(ctx context.Context, org, repo string, number int, comment string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteComment(org, repo string, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteCommentWithContext(ctx context.Context, org, repo string, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditComment(org, repo string, id int, comment string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditCommentWithContext(ctx context.Context, org, repo string, id int, comment string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateCommentReaction(org, repo string, id int, reaction string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteStaleComments(org, repo string, number int, comments []github.IssueComment, isStale func(github.IssueComment) bool) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteStaleCommentsWithContext(ctx context.Context, org, repo string, number int, comments []github.IssueComment, isStale func(github.IssueComment) bool) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) IsMember(org, user string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetOrg(name string) (*github.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditOrg(name string, config github.Organization) (*github.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListOrgInvitations(org string) ([]github.OrgInvitation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListOrgMembers(org, role string) ([]github.TeamMember, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) HasPermission(org, repo, user string, roles ...string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetUserPermission(org, repo, user string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateOrgMembership(org, user string, admin bool) (*github.OrgMembership, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveOrgMembership(org, user string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateTeam(org string, team github.Team) (*github.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditTeam(org string, t github.Team) (*github.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteTeam(org string, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteTeamBySlug(org, teamSlug string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeams(org string) ([]github.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateTeamMembership(org string, id int, user string, maintainer bool) (*github.TeamMembership, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateTeamMembershipBySlug(org, teamSlug, user string, maintainer bool) (*github.TeamMembership, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveTeamMembership(org string, id int, user string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveTeamMembershipBySlug(org, teamSlug, user string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeamMembers(org string, id int, role string) ([]github.TeamMember, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeamMembersBySlug(org, teamSlug, role string) ([]github.TeamMember, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeamRepos(org string, id int) ([]github.Repo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeamReposBySlug(org, teamSlug string) ([]github.Repo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateTeamRepo(id int, org, repo string, permission github.TeamPermission) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) UpdateTeamRepoBySlug(org, teamSlug, repo string, permission github.TeamPermission) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveTeamRepo(id int, org, repo string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) RemoveTeamRepoBySlug(org, teamSlug, repo string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeamInvitations(org string, id int) ([]github.OrgInvitation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListTeamInvitationsBySlug(org, teamSlug string) ([]github.OrgInvitation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) TeamHasMember(org string, teamID int, memberLogin string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) TeamBySlugHasMember(org string, teamSlug string, memberLogin string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetTeamBySlug(slug string, org string) (*github.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetRepoProjects(owner, repo string) ([]github.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetOrgProjects(org string) ([]github.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetProjectColumns(org string, projectID int) ([]github.ProjectColumn, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateProjectCard(org string, columnID int, projectCard github.ProjectCard) (*github.ProjectCard, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetColumnProjectCards(org string, columnID int) ([]github.ProjectCard, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetColumnProjectCard(org string, columnID int, issueURL string) (*github.ProjectCard, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) MoveProjectCard(org string, projectCardID int, newColumnID int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteProjectCard(org string, projectCardID int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ClearMilestone(org, repo string, num int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) SetMilestone(org, repo string, issueNum, milestoneNum int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListMilestones(org, repo string) ([]github.Milestone, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) BotUserChecker() (func(candidate string) bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) BotUserCheckerWithContext(ctx context.Context) (func(candidate string) bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) Email() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListOrgHooks(org string) ([]github.Hook, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListRepoHooks(org, repo string) ([]github.Hook, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditRepoHook(org, repo string, id int, req github.HookRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) EditOrgHook(org string, id int, req github.HookRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateOrgHook(org string, req github.HookRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) CreateRepoHook(org, repo string, req github.HookRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteOrgHook(org string, id int, req github.HookRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteRepoHook(org, repo string, id int, req github.HookRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListCurrentUserRepoInvitations() ([]github.UserRepoInvitation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) AcceptUserRepoInvitation(invitationID int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListCurrentUserOrgInvitations() ([]github.UserOrgInvitation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) AcceptUserOrgInvitation(org string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListAppInstallations() ([]github.AppInstallation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) IsAppInstalled(org, repo string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) UsesAppAuth() bool {
	//TODO implement me
	panic("implement me")
}

func (c *client) ListAppInstallationsForOrg(org string) ([]github.AppInstallation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetApp() (*github.App, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetFailedActionRunsByHeadBranch(org, repo, branchName, headSHA string) ([]github.WorkflowRun, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) Throttle(hourlyTokens, burst int, org ...string) error {
	//TODO implement me
	panic("implement me")
}

// QueryWithGitHubAppsSupport runs a GraphQL query using shurcooL/githubql's client.
func (c *client) QueryWithGitHubAppsSupport(ctx context.Context, q interface{}, vars map[string]interface{}, org string) error {
	// Don't log query here because Query is typically called multiple times to get all pages.
	// Instead log once per search and include total search cost.
	return c.gqlc.QueryWithGitHubAppsSupport(ctx, q, vars, org)
}

func (c *client) MutateWithGitHubAppsSupport(ctx context.Context, m interface{}, input githubql.Input, vars map[string]interface{}, org string) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) SetMax404Retries(i int) {
	//TODO implement me
	panic("implement me")
}

func (c *client) WithFields(fields logrus.Fields) github.Client {
	//TODO implement me
	panic("implement me")
}

func (c *client) ForPlugin(plugin string) github.Client {
	//TODO implement me
	panic("implement me")
}

func (c *client) ForSubcomponent(subcomponent string) github.Client {
	//TODO implement me
	panic("implement me")
}

func (c *client) Used() bool {
	//TODO implement me
	panic("implement me")
}

func (c *client) TriggerGitHubWorkflow(org, repo string, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) TriggerFailedGitHubWorkflow(org, repo string, id int) error {
	//TODO implement me
	panic("implement me")
}

// This should only be called once when the client is created.
func (c *client) wrapThrottler() {
	c.throttle.http = c.client
	c.throttle.graph = c.gqlc
	c.client = &c.throttle
	c.gqlc = &c.throttle
}

func (c *client) userAgent() string {
	if c.identifier != "" {
		return version.UserAgentWithIdentifier(c.identifier)
	}
	return version.UserAgent()
}

// BotUser returns the user data of the authenticated identity.
//
// See https://developer.github.com/v3/users/#get-the-authenticated-user
func (c *client) BotUser() (*github.UserData, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.userData == nil {
		if err := c.getUserData(context.Background()); err != nil {
			return nil, fmt.Errorf("fetching bot name from GitHub: %w", err)
		}
	}
	return c.userData, nil
}

// Not thread-safe - callers need to hold c.mut.
func (c *client) getUserData(ctx context.Context) error {
	if c.delegate.usesAppsAuth {
		resp, err := c.GetAppWithContext(ctx)
		if err != nil {
			return err
		}
		c.userData = &github.UserData{
			Name:  resp.Name,
			Login: resp.Slug,
			Email: fmt.Sprintf("%s@users.noreply.github.com", resp.Slug),
		}
		return nil
	}
	c.log("User")
	var u github.User
	_, err := c.requestWithContext(ctx, &request{
		method:    http.MethodGet,
		path:      "/user",
		exitCodes: []int{200},
	}, &u)
	if err != nil {
		return err
	}
	c.userData = &github.UserData{
		Name:  u.Name,
		Login: u.Login,
		Email: u.Email,
	}
	// email needs to be publicly accessible via the profile
	// of the current account. Read below for more info
	// https://developer.github.com/v3/users/#get-a-single-user

	// record information for the user
	authHeaderHash := fmt.Sprintf("%x", sha256.Sum256([]byte(c.authHeader()))) // use %x to make this a utf-8 string for use as a label
	userInfo.With(prometheus.Labels{"token_hash": authHeaderHash, "login": c.userData.Login, "email": c.userData.Email}).Set(1)
	return nil
}

func (c *client) GetAppWithContext(ctx context.Context) (*github.App, error) {
	durationLogger := c.log("App")
	defer durationLogger()

	var app github.App
	if _, err := c.requestWithContext(ctx, &request{
		method:    http.MethodGet,
		path:      "/app",
		exitCodes: []int{200},
	}, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

func (c *client) requestWithContext(ctx context.Context, r *request, ret interface{}) (int, error) {
	statusCode, b, err := c.requestRawWithContext(ctx, r)
	if err != nil {
		return statusCode, err
	}
	if ret != nil {
		if err := json.Unmarshal(b, ret); err != nil {
			return statusCode, err
		}
	}
	return statusCode, nil
}

func (c *client) requestRawWithContext(ctx context.Context, r *request) (int, []byte, error) {
	if c.fake || (c.dry && r.method != http.MethodGet) {
		return r.exitCodes[0], nil, nil
	}
	resp, err := c.requestRetryWithContext(ctx, r.method, r.path, r.accept, r.org, r.requestBody)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	var okCode bool
	for _, code := range r.exitCodes {
		if code == resp.StatusCode {
			okCode = true
			break
		}
	}
	if !okCode {
		clientError := unmarshalClientError(b)
		err = requestError{
			StatusCode:  resp.StatusCode,
			ClientError: clientError,
			ErrorString: fmt.Sprintf("status code %d not one of %v, body: %s", resp.StatusCode, r.exitCodes, string(b)),
		}
	}
	return resp.StatusCode, b, err
}

func (c *client) requestRetryWithContext(ctx context.Context, method, path, accept, org string, body interface{}) (*http.Response, error) {
	var hostIndex int
	var resp *http.Response
	var err error
	backoff := c.initialDelay
	for retries := 0; retries < c.maxRetries; retries++ {
		if retries > 0 && resp != nil {
			resp.Body.Close()
		}
		resp, err = c.doRequest(ctx, method, c.bases[hostIndex]+path, accept, org, body)
		if err == nil {
			if resp.StatusCode == 404 && retries < c.max404Retries {
				// Retry 404s a couple times. Sometimes GitHub is inconsistent in
				// the sense that they send us an event such as "PR opened" but an
				// immediate request to GET the PR returns 404. We don't want to
				// retry more than a couple times in this case, because a 404 may
				// be caused by a bad API call and we'll just burn through API
				// tokens.
				c.logger.WithField("backoff", backoff.String()).Debug("Retrying 404")
				c.time.Sleep(backoff)
				backoff *= 2
			} else if resp.StatusCode == 403 {
				if resp.Header.Get("X-RateLimit-Remaining") == "0" {
					// If we are out of API tokens, sleep first. The X-RateLimit-Reset
					// header tells us the time at which we can request again.
					var t int
					if t, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Reset")); err == nil {
						// Sleep an extra second plus how long GitHub wants us to
						// sleep. If it's going to take too long, then break.
						sleepTime := c.time.Until(time.Unix(int64(t), 0)) + time.Second
						if sleepTime < c.maxSleepTime {
							c.logger.WithField("backoff", sleepTime.String()).WithField("path", path).Debug("Retrying after token budget reset")
							c.time.Sleep(sleepTime)
						} else {
							err = fmt.Errorf("sleep time for token reset exceeds max sleep time (%v > %v)", sleepTime, c.maxSleepTime)
							resp.Body.Close()
							break
						}
					} else {
						err = fmt.Errorf("failed to parse rate limit reset unix time %q: %w", resp.Header.Get("X-RateLimit-Reset"), err)
						resp.Body.Close()
						break
					}
				} else if rawTime := resp.Header.Get("Retry-After"); rawTime != "" && rawTime != "0" {
					// If we are getting abuse rate limited, we need to wait or
					// else we risk continuing to make the situation worse
					var t int
					if t, err = strconv.Atoi(rawTime); err == nil {
						// Sleep an extra second plus how long GitHub wants us to
						// sleep. If it's going to take too long, then break.
						sleepTime := time.Duration(t+1) * time.Second
						if sleepTime < c.maxSleepTime {
							c.logger.WithField("backoff", sleepTime.String()).WithField("path", path).Debug("Retrying after abuse ratelimit reset")
							c.time.Sleep(sleepTime)
						} else {
							err = fmt.Errorf("sleep time for abuse rate limit exceeds max sleep time (%v > %v)", sleepTime, c.maxSleepTime)
							resp.Body.Close()
							break
						}
					} else {
						err = fmt.Errorf("failed to parse abuse rate limit wait time %q: %w", rawTime, err)
						resp.Body.Close()
						break
					}
				} else {
					acceptedScopes := resp.Header.Get("X-Accepted-OAuth-Scopes")
					authorizedScopes := resp.Header.Get("X-OAuth-Scopes")
					if authorizedScopes == "" {
						authorizedScopes = "no"
					}

					want := sets.New[string]()
					for _, acceptedScope := range strings.Split(acceptedScopes, ",") {
						want.Insert(strings.TrimSpace(acceptedScope))
					}
					var got []string
					for _, authorizedScope := range strings.Split(authorizedScopes, ",") {
						got = append(got, strings.TrimSpace(authorizedScope))
					}
					if acceptedScopes != "" && !want.HasAny(got...) {
						err = fmt.Errorf("the account is using %s oauth scopes, please make sure you are using at least one of the following oauth scopes: %s", authorizedScopes, acceptedScopes)
					} else {
						body, _ := io.ReadAll(resp.Body)
						err = fmt.Errorf("the GitHub API request returns a 403 error: %s", string(body))
					}
					resp.Body.Close()
					break
				}
			} else if resp.StatusCode < 500 {
				// Normal, happy case.
				break
			} else {
				// Retry 500 after a break.
				c.logger.WithField("backoff", backoff.String()).Debug("Retrying 5XX")
				c.time.Sleep(backoff)
				backoff *= 2
			}
		} else if errors.Is(err, &appsAuthError{}) {
			c.logger.WithError(err).Error("Stopping retry due to appsAuthError")
			return resp, err
		} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return resp, err
		} else {
			// Connection problem. Try a different host.
			oldHostIndex := hostIndex
			hostIndex = (hostIndex + 1) % len(c.bases)
			c.logger.WithFields(logrus.Fields{
				"err":          err,
				"backoff":      backoff.String(),
				"old-endpoint": c.bases[oldHostIndex],
				"new-endpoint": c.bases[hostIndex],
			}).Debug("Retrying request due to connection problem")
			c.time.Sleep(backoff)
			backoff *= 2
		}
	}
	return resp, err
}

func (c *client) doRequest(ctx context.Context, method, path, accept, org string, body interface{}) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		b = c.censor(b)
		buf = bytes.NewBuffer(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, path, buf)
	if err != nil {
		return nil, fmt.Errorf("failed creating new request: %w", err)
	}
	// We do not make use of the Set() method to set this header because
	// the header name `X-GitHub-Api-Version` is non-canonical in nature.
	//
	// See https://pkg.go.dev/net/http#Header.Set for more info.
	req.Header["X-GitHub-Api-Version"] = []string{githubApiVersion}
	c.logger.Debugf("Using GitHub REST API Version: %s", githubApiVersion)
	if header := c.authHeader(); len(header) > 0 {
		req.Header.Set("Authorization", header)
	}
	if accept == acceptNone {
		req.Header.Add("Accept", "application/vnd.github.v3+json")
	} else {
		req.Header.Add("Accept", accept)
	}
	if userAgent := c.userAgent(); userAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}
	if org != "" {
		req = req.WithContext(context.WithValue(req.Context(), githubOrgHeaderKey, org))
	}
	// Disable keep-alive so that we don't get flakes when GitHub closes the
	// connection prematurely.
	// https://go-review.googlesource.com/#/c/3210/ fixed it for GET, but not
	// for POST.
	req.Close = true

	c.logger.WithField("curl", toCurl(req)).Trace("Executing http request")
	return c.client.Do(req)
}

func (c *client) authHeader() string {
	if c.getToken == nil {
		return ""
	}
	token := c.getToken()
	if len(token) == 0 {
		return ""
	}
	return fmt.Sprintf("Bearer %s", token)
}

func (c *client) log(methodName string, args ...interface{}) (logDuration func()) {
	c.mutUsed.Lock()
	c.used = true
	c.mutUsed.Unlock()

	if c.logger == nil {
		return func() {}
	}
	var as []string
	for _, arg := range args {
		as = append(as, fmt.Sprintf("%v", arg))
	}
	start := time.Now()
	c.logger.Infof("%s(%s)", methodName, strings.Join(as, ", "))
	return func() {
		c.logger.WithField("duration", time.Since(start).String()).Debugf("%s(%s) finished", methodName, strings.Join(as, ", "))
	}
}

// Make a request with retries. If ret is not nil, unmarshal the response body
// into it. Returns an error if the exit code is not one of the provided codes.
func (c *client) request(r *request, ret interface{}) (int, error) {
	return c.requestWithContext(context.Background(), r, ret)
}

// readPaginatedResults iterates over all objects in the paginated result indicated by the given url.
//
// newObj() should return a new slice of the expected type
// accumulate() should accept that populated slice for each page of results.
//
// Returns an error any call to GitHub or object marshalling fails.
func (c *client) readPaginatedResults(path, accept, org string, newObj func() interface{}, accumulate func(interface{})) error {
	return c.readPaginatedResultsWithContext(context.Background(), path, accept, org, newObj, accumulate)
}

func (c *client) readPaginatedResultsWithContext(ctx context.Context, path, accept, org string, newObj func() interface{}, accumulate func(interface{})) error {
	values := url.Values{
		"per_page": []string{"100"},
	}
	return c.readPaginatedResultsWithValuesWithContext(ctx, path, values, accept, org, newObj, accumulate)
}

// readPaginatedResultsWithValues is an override that allows control over the query string.
/*func (c *client) readPaginatedResultsWithValues(path string, values url.Values, accept, org string, newObj func() interface{}, accumulate func(interface{})) error {
	return c.readPaginatedResultsWithValuesWithContext(context.Background(), path, values, accept, org, newObj, accumulate)
}*/

func (c *client) readPaginatedResultsWithValuesWithContext(ctx context.Context, path string, values url.Values, accept, org string, newObj func() interface{}, accumulate func(interface{})) error {
	pagedPath := path
	if len(values) > 0 {
		pagedPath += "?" + values.Encode()
	}
	for {
		resp, err := c.requestRetryWithContext(ctx, http.MethodGet, pagedPath, accept, org, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("return code not 2XX: %s", resp.Status)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		obj := newObj()
		if err := json.Unmarshal(b, obj); err != nil {
			return err
		}

		accumulate(obj)

		link := parseLinks(resp.Header.Get("Link"))["next"]
		if link == "" {
			break
		}

		// Example for github.com:
		// * c.bases[0]: api.github.com
		// * initial call: api.github.com/repos/kubernetes/kubernetes/pulls?per_page=100
		// * next: api.github.com/repositories/22/pulls?per_page=100&page=2
		// * in this case prefix will be empty and we're just calling the path returned by next
		// Example for github enterprise:
		// * c.bases[0]: <ghe-url>/api/v3
		// * initial call: <ghe-url>/api/v3/repos/kubernetes/kubernetes/pulls?per_page=100
		// * next: <ghe-url>/api/v3/repositories/22/pulls?per_page=100&page=2
		// * in this case prefix will be "/api/v3" and we will strip the prefix. If we don't do that,
		//   the next call will go to <ghe-url>/api/v3/api/v3/repositories/22/pulls?per_page=100&page=2
		prefix := strings.TrimSuffix(resp.Request.URL.RequestURI(), pagedPath)

		u, err := url.Parse(link)
		if err != nil {
			return fmt.Errorf("failed to parse 'next' link: %w", err)
		}
		pagedPath = strings.TrimPrefix(u.RequestURI(), prefix)
	}
	return nil
}

func (c *client) getAppInstallationToken(installationId int64) (*github.AppInstallationToken, error) {
	durationLogger := c.log("AppInstallationToken")
	defer durationLogger()

	if c.dry {
		return nil, fmt.Errorf("not requesting GitHub App access_token in dry-run mode")
	}

	var token github.AppInstallationToken
	if _, err := c.request(&request{
		method:    http.MethodPost,
		path:      fmt.Sprintf("/app/installations/%d/access_tokens", installationId),
		exitCodes: []int{201},
	}, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// GetPullRequests get all open pull requests for a repo.
//
// See https://developer.github.com/v3/pulls/#list-pull-requests
func (c *client) GetPullRequests(org, repo string) ([]github.PullRequest, error) {
	c.log("GetPullRequests", org, repo)
	var prs []github.PullRequest
	if c.fake {
		return prs, nil
	}
	path := fmt.Sprintf("/repos/%s/%s/pulls", org, repo)
	err := c.readPaginatedResults(
		path,
		// allow the description and draft fields
		// https://developer.github.com/changes/2018-02-22-label-description-search-preview/
		// https://developer.github.com/changes/2019-02-14-draft-pull-requests/
		"application/vnd.github.symmetra-preview+json, application/vnd.github.shadow-cat-preview",
		org,
		func() interface{} {
			return &[]github.PullRequest{}
		},
		func(obj interface{}) {
			prs = append(prs, *(obj.(*[]github.PullRequest))...)
		},
	)
	if err != nil {
		return nil, err
	}
	return prs, err
}

// GetPullRequest gets a pull request.
//
// See https://developer.github.com/v3/pulls/#get-a-single-pull-request
func (c *client) GetPullRequest(org, repo string, number int) (*github.PullRequest, error) {
	durationLogger := c.log("GetPullRequest", org, repo, number)
	defer durationLogger()

	var pr github.PullRequest
	_, err := c.request(&request{
		// allow the description and draft fields
		// https://developer.github.com/changes/2018-02-22-label-description-search-preview/
		// https://developer.github.com/changes/2019-02-14-draft-pull-requests/
		accept:    "application/vnd.github.symmetra-preview+json, application/vnd.github.shadow-cat-preview",
		method:    http.MethodGet,
		path:      fmt.Sprintf("/repos/%s/%s/pulls/%d", org, repo, number),
		org:       org,
		exitCodes: []int{200},
	}, &pr)
	return &pr, err
}

type request struct {
	method      string
	path        string
	accept      string
	org         string
	requestBody interface{}
	exitCodes   []int
}

type requestError struct {
	StatusCode  int
	ClientError error
	ErrorString string
}

func (r requestError) Error() string {
	return r.ErrorString
}

func (r requestError) ErrorMessages() []string {
	clientErr, isClientError := r.ClientError.(github.ClientError)
	if isClientError {
		errors := []string{}
		for _, subErr := range clientErr.Errors {
			errors = append(errors, subErr.Message)
		}
		return errors
	}
	alternativeClientErr, isAlternativeClientError := r.ClientError.(github.AlternativeClientError)
	if isAlternativeClientError {
		return alternativeClientErr.Errors
	}
	return []string{}
}

// Interface for how prow interacts with the graphql client, which we may throttle.
type gqlClient interface {
	QueryWithGitHubAppsSupport(ctx context.Context, q interface{}, vars map[string]interface{}, org string) error
	MutateWithGitHubAppsSupport(ctx context.Context, m interface{}, input githubql.Input, vars map[string]interface{}, org string) error
	forUserAgent(userAgent string) gqlClient
}

type graphQLGitHubAppsAuthClientWrapper struct {
	*githubql.Client
	userAgent string
}

func (c *graphQLGitHubAppsAuthClientWrapper) QueryWithGitHubAppsSupport(ctx context.Context, q interface{}, vars map[string]interface{}, org string) error {
	ctx = context.WithValue(ctx, githubOrgHeaderKey, org)
	ctx = context.WithValue(ctx, userAgentContextKey, c.userAgent)
	return c.Client.Query(ctx, q, vars)
}

func (c *graphQLGitHubAppsAuthClientWrapper) MutateWithGitHubAppsSupport(ctx context.Context, m interface{}, input githubql.Input, vars map[string]interface{}, org string) error {
	ctx = context.WithValue(ctx, githubOrgHeaderKey, org)
	ctx = context.WithValue(ctx, userAgentContextKey, c.userAgent)
	return c.Client.Mutate(ctx, m, input, vars)
}

func (c *graphQLGitHubAppsAuthClientWrapper) forUserAgent(userAgent string) gqlClient {
	return &graphQLGitHubAppsAuthClientWrapper{
		Client:    c.Client,
		userAgent: userAgent,
	}
}

// userInfo provides the 'github_user_info' vector that is indexed
// by the user's information.
var userInfo = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "github_user_info",
		Help: "Metadata about a user, tied to their token hash.",
	},
	[]string{"token_hash", "login", "email"},
)
