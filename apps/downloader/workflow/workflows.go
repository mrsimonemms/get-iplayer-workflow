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

package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func DownloadBBCProgramme(ctx workflow.Context, download Download) (any, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Downloading BBC programme", "pid", download.ProgrammeID)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
	})

	logger.Debug("Downloading programme")
	var downloadByPIDResult *DownloadByPIDResult
	if err := workflow.ExecuteActivity(ctx, DownloadByPID, download).Get(ctx, &downloadByPIDResult); err != nil {
		logger.Error("Error downloading programme", "pid", download.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error downloading programme: %w", err)
	}

	// Invoke the child workflows in parallel
	futures := map[workflow.Context]workflow.ChildWorkflowFuture{}
	for i, f := range downloadByPIDResult.Files {
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowTaskTimeout: time.Hour,
			WorkflowID:          fmt.Sprintf("%s_parse_%d", workflow.GetInfo(ctx).WorkflowExecution.ID, i),
		})

		file := DownloadedProgramme{
			ProgrammeID: downloadByPIDResult.ProgrammeID,
			SavePath:    downloadByPIDResult.SavePath,
			File:        f,
		}
		futures[childCtx] = workflow.ExecuteChildWorkflow(childCtx, ParseDownloadedProgramme, file)
	}

	// Now the child workflows are running, wait for the results
	// for ctx, workflow := range futures {
	// 	var node *providers.NodeResult

	// 	if err := workflow.Get(ctx, &node); err != nil {
	// 		logger.Error("Error provisioning nodes", "error", err)
	// 		return nil, fmt.Errorf("error provisioning nodes: %w", err)
	// 	}

	// 	project.Nodes = append(project.Nodes, node)
	// }

	return "hello", nil
}

func ParseDownloadedProgramme(ctx workflow.Context, programme DownloadedProgramme) (any, error) {
	// logger := workflow.GetLogger(ctx)
	// logger.Info("Downloading BBC programme", "pid", download.ProgrammeID)

	// ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
	// 	StartToCloseTimeout: time.Hour,
	// 	RetryPolicy: &temporal.RetryPolicy{
	// 		InitialInterval:    time.Second,
	// 		BackoffCoefficient: 2.0,
	// 		MaximumInterval:    time.Minute * 5,
	// 		MaximumAttempts:    3,
	// 	},
	// })

	// Upload to S3 bucket

	// Convert file to Plex format

	// Update the audio headers

	// Upload to target location

	return nil, nil
}
