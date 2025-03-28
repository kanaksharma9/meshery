// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package environments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/environments"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environments",
	Long:  `Create a new environments by providing the name and description of the environment`,
	Example: `
// Create a new environment
mesheryctl environment create --orgID [orgID] --name [name] --description [description]
// Documentation for environment can be found at:
https://docs.meshery.io/concepts/logical/environments
`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Check if all three flags are set
		orgIDFlag, _ := cmd.Flags().GetString("orgID")
		nameFlag, _ := cmd.Flags().GetString("name")
		descriptionFlag, _ := cmd.Flags().GetString("description")

		if orgIDFlag == "" || nameFlag == "" || descriptionFlag == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("Please provide a --orgID, --name, and --description flag"))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/environments", baseUrl)

		if name == "" || description == "" {
			return utils.ErrInvalidArgument(errors.New("name is required"))
		}

		payload := &environments.EnvironmentPayload{
			Name:        name,
			Description: description,
			OrgID:       orgID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		req, err := utils.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusOK {
			utils.Log.Info("environment created")
			return nil
		}
		utils.Log.Info("Error creating environment")
		return nil
	},
}
