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

	return nil, nil
}
