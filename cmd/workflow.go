/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"github.com/mrsimonemms/get-iplayer-workflow/pkg/temporal"
	"github.com/mrsimonemms/get-iplayer-workflow/pkg/workflows"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/worker"
)

func addWorkflowCmd() *cobra.Command {
	// workflowCmd represents the workflow command
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Run the workflow",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := temporal.New()
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to create Temporal client")
			}

			defer client.Close()

			w := workflows.NewGetIplayerWorkflow(client)

			if err := w.Run(worker.InterruptCh()); err != nil {
				log.Fatal().Err(err).Msg("Failed to start workflow run")
			}
		},
	}

	return workflowCmd
}
