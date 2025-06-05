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

	"github.com/minio/minio-go/v7"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func DownloadBBCProgramme(ctx workflow.Context, download Download) (*DownloadBBCProgrammeResult, error) {
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

	var a *activities

	logger.Debug("Downloading programme")
	var downloadByPIDResult *DownloadByPIDResult
	if err := workflow.ExecuteActivity(ctx, a.DownloadByPID, download).Get(ctx, &downloadByPIDResult); err != nil {
		logger.Error("Error downloading programme", "pid", download.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error downloading programme: %w", err)
	}

	// Invoke the child workflows in parallel
	parentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	futures := map[workflow.Context]workflow.ChildWorkflowFuture{}
	for i, f := range downloadByPIDResult.Files {
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowTaskTimeout: time.Hour,
			WorkflowID:          fmt.Sprintf("%s_parse_%d", parentWorkflowID, i),
		})

		file := DownloadedProgramme{
			ProgrammeID: downloadByPIDResult.ProgrammeID,
			SavePath:    downloadByPIDResult.SavePath,
			File:        f,
		}
		futures[childCtx] = workflow.ExecuteChildWorkflow(childCtx, ParseDownloadedProgramme, file, parentWorkflowID)
	}

	// Now the child workflows are running, wait for the results
	for ctx, workflow := range futures {
		var result *ParseDownloadedProgrammeResult
		if err := workflow.Get(ctx, &result); err != nil {
			logger.Error("Error parsing download", "error", err)
			return nil, fmt.Errorf("error parsing download: %w", err)
		}

		fmt.Printf("%+v\n", result)
	}

	return &DownloadBBCProgrammeResult{}, nil
}

func ParseDownloadedProgramme(ctx workflow.Context, programme DownloadedProgramme, parentWorkflowID string) (*ParseDownloadedProgrammeResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Parsing downloaded BBC programme", "pid", programme.ProgrammeID)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
	})

	var a *activities

	// Get file name
	logger.Debug("Generating programme name")
	var programmeNameResult *ProgrammeNameResult
	if err := workflow.ExecuteActivity(ctx, a.GenerateProgrammeName, programme, parentWorkflowID).Get(ctx, &programmeNameResult); err != nil {
		logger.Error("Error generating programme name", "pid", programme.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error generating programme name: %w", err)
	}
	programme.TargetName = programmeNameResult.Name

	// Upload to S3 bucket
	logger.Debug("Uploading programme to S3 bucket")
	var uploadedResult *minio.UploadInfo
	if err := workflow.ExecuteActivity(ctx, a.UploadFileToS3Bucket, programme, parentWorkflowID).Get(ctx, &uploadedResult); err != nil {
		logger.Error("Error uploading programme", "pid", programme.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error uploading programme: %w", err)
	}

	// Update the audio headers
	if programmeNameResult.API.Programme.MediaType == BBCProgrammeAPIProgrammeMediaTypeAudio {
	}

	// Upload to target location

	return &ParseDownloadedProgrammeResult{
		ProgrammeID: programme.ProgrammeID,
		Bucket:      uploadedResult,
	}, nil
}
