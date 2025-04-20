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

package workflows

import (
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func DownloadItemFromIplayerWorkflow(ctx workflow.Context, cfg DownloadItemConfig) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Download item from iPlayer")
}

func DownloadURLFromIplayerWorkflow(ctx workflow.Context, cfg DownloadURLConfig) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting new parse iPlayer URL workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	})

	logger.Debug("Parse the URL to get the programme IDs")
}

func NewGetIplayerWorkflow(client client.Client) worker.Worker {
	w := worker.New(client, GetIplayerWorkflowName, worker.Options{})

	// Register the workflows
	w.RegisterWorkflow(DownloadURLFromIplayerWorkflow)

	return w
}
