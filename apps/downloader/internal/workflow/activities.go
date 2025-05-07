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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/mrsimonemms/get-iplayer-workflow/apps/downloader/internal/config"
	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/activity"
)

type activities struct {
	cfg   *config.Config
	nc    *nats.Conn
	minio *minio.Client
}

func (a *activities) DownloadByPID(ctx context.Context, download Download) (*DownloadByPIDResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Downloading programme by PID", "pid", download.ProgrammeID)

	workflowID := activity.GetInfo(ctx).WorkflowExecution.ID

	savePath := path.Join(a.cfg.OutputDir, workflowID)

	args := []string{
		"--nocopyright",
		"--subdir",
		"--whitespace",
		"--pid-recursive",
		"--force",
		fmt.Sprintf("--pid=%s", download.ProgrammeID),
		fmt.Sprintf("--output=%s", savePath),
	}

	logger.Debug("Command", "args", strings.Join(args, " "))

	so := &streamOutput{
		logger:     logger,
		nc:         a.nc,
		workflowID: workflowID,
	}
	cmd := exec.CommandContext(ctx, "get_iplayer", args...)
	cmd.Stdout = so
	cmd.Stderr = so
	if err := cmd.Start(); err != nil {
		logger.Error("Error downloading with get_iplayer", "pid", download.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error downloading with get_iplayer: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		logger.Error("Error executing get_iplayer", "pid", download.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error executiing get_iplayer: %w", err)
	}

	logger.Info("Programme downloaded", "pid", download.ProgrammeID)
	_, _ = so.Write([]byte("Programme downloaded"))

	files := make([]string, 0)

	if err := filepath.Walk(savePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	}); err != nil {
		logger.Error("Error listing files", "pid", download.ProgrammeID, "error", err)
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	return &DownloadByPIDResult{
		ProgrammeID: download.ProgrammeID,
		SavePath:    savePath,
		Files:       files,
	}, nil
}

func (a *activities) GenerateProgrammeName(ctx context.Context, programme DownloadedProgramme, _ string) (*ProgrammeNameResult, error) {
	logger := activity.GetLogger(ctx)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("https://www.bbc.co.uk/programmes/%s.json", programme.ProgrammeID)
	logger.Debug("Making call to BBC API for programme name", "url", url)

	res, err := client.Get(url)
	if err != nil {
		logger.Error("Error making HTTP request", "error", err)
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	if res.StatusCode != 200 {
		logger.Error("Unknown programme ID", "pid", programme.ProgrammeID)
		return nil, fmt.Errorf("unknown programme id: %s", programme.ProgrammeID)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error parsing HTTP data", "error", err)
		return nil, fmt.Errorf("error parsing http data: %w", err)
	}

	data := BBCProgrammeAPI{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error("Error converting to programme result", "error", err)
		return nil, fmt.Errorf("error converting to programme result")
	}

	return &ProgrammeNameResult{
		Name: data.GetFileName(filepath.Ext(programme.File)),
	}, nil
}

func (a *activities) UploadFileToS3Bucket(ctx context.Context, programme DownloadedProgramme, parentWorkflowID string) (*minio.UploadInfo, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Detecting mime type of file", "file", programme.File)
	mtype, err := mimetype.DetectFile(programme.File)
	if err != nil {
		logger.Error("Error detecting mime type of file", "file", programme.File)
		return nil, fmt.Errorf("error detecting mime type of file: %w", err)
	}

	logger.Info("Uploading file to S3 bucket", "mimetype", mtype.String(), "file", programme.File)
	_ = a.nc.Publish(fmt.Sprintf("log.%s", parentWorkflowID), fmt.Appendf(nil, "Uploading: %s", programme.TargetName))
	info, err := a.minio.FPutObject(ctx, a.cfg.S3.Bucket, programme.TargetName, programme.File, minio.PutObjectOptions{
		ContentType: mtype.String(),
	})
	if err != nil {
		logger.Error("Error uploading file to bucket", "error", err)
		return nil, fmt.Errorf("error uploading file to bucket: %w", err)
	}
	_ = a.nc.Publish(fmt.Sprintf("log.%s", parentWorkflowID), fmt.Appendf(nil, "Uploaded: %s", programme.TargetName))

	return &info, nil
}

func NewActivities(nc *nats.Conn, cfg *config.Config, m *minio.Client) *activities {
	return &activities{
		nc:    nc,
		cfg:   cfg,
		minio: m,
	}
}
