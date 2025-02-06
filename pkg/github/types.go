package github

import (
	"encoding/json"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/prow/pkg/github"
)

// User is a GitHub user account.
/*type User struct {
	Login       string          `json:"login"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	ID          int             `json:"id"`
	HTMLURL     string          `json:"html_url"`
	Permissions RepoPermissions `json:"permissions"`
	Type        string          `json:"type"`
}

// RepoPermissions describes which permission level an entity has in a
// repo. At most one of the booleans here should be true.
type RepoPermissions struct {
	// Pull is equivalent to "Read" permissions in the web UI
	Pull   bool `json:"pull"`
	Triage bool `json:"triage"`
	// Push is equivalent to "Edit" permissions in the web UI
	Push     bool `json:"push"`
	Maintain bool `json:"maintain"`
	Admin    bool `json:"admin"`
}

// AppInstallation represents a GitHub Apps installation.
type AppInstallation struct {
	ID                  int64                   `json:"id,omitempty"`
	AppSlug             string                  `json:"app_slug,omitempty"`
	NodeID              string                  `json:"node_id,omitempty"`
	AppID               int64                   `json:"app_id,omitempty"`
	TargetID            int64                   `json:"target_id,omitempty"`
	Account             User                    `json:"account,omitempty"`
	AccessTokensURL     string                  `json:"access_tokens_url,omitempty"`
	RepositoriesURL     string                  `json:"repositories_url,omitempty"`
	HTMLURL             string                  `json:"html_url,omitempty"`
	TargetType          string                  `json:"target_type,omitempty"`
	SingleFileName      string                  `json:"single_file_name,omitempty"`
	RepositorySelection string                  `json:"repository_selection,omitempty"`
	Events              []string                `json:"events,omitempty"`
	Permissions         InstallationPermissions `json:"permissions,omitempty"`
	CreatedAt           string                  `json:"created_at,omitempty"`
	UpdatedAt           string                  `json:"updated_at,omitempty"`
}*/

func unmarshalClientError(b []byte) error {
	var errors []error
	clientError := github.ClientError{}
	err := json.Unmarshal(b, &clientError)
	if err == nil {
		return clientError
	}
	errors = append(errors, err)
	alternativeClientError := github.AlternativeClientError{}
	err = json.Unmarshal(b, &alternativeClientError)
	if err == nil {
		return alternativeClientError
	}
	errors = append(errors, err)
	return utilerrors.NewAggregate(errors)
}
