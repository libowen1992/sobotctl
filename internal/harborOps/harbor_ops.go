package harborOps

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/mittwald/goharbor-client/v5/apiv2"
	"sobotctl/global"
)

type HarborOps struct {
	Client *apiv2.RESTClient
}

func NewHarborOps() (*HarborOps, error) {
	harborC, err := apiv2.NewRESTClientForHost(
		fmt.Sprintf("%s/api", global.HarborS.URL), global.HarborS.User, global.HarborS.Password, nil)
	if err != nil {
		return nil, err
	}
	return &HarborOps{
		Client: harborC,
	}, nil
}

type HarborProject struct {
	// The creation time of the project.
	// Format: date-time
	CreationTime strfmt.DateTime `json:"creation_time,omitempty"`

	// The name of the project.
	Name string `json:"name,omitempty"`

	// The number of the repositories under this project.
	RepoCount int64 `json:"repo_count"`

	// The update time of the project.
	// Format: date-time
	UpdateTime strfmt.DateTime `json:"update_time,omitempty"`

	// The public status of the project. The valid values are "true", "false".
	Public string `json:"public,omitempty"`
}
