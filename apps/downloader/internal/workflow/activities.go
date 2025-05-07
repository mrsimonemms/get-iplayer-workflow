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
	"fmt"
	"io/fs"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/mrsimonemms/get-iplayer-workflow/apps/downloader/internal/config"
	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/activity"
)

type activities struct {
	nc *nats.Conn
}

func (a *activities) DownloadByPID(ctx context.Context, download Download) (*DownloadByPIDResult, error) {
	logger := activity.GetLogger(ctx)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	logger.Info("Downloading programme by PID", "pid", download.ProgrammeID)

	workflowID := activity.GetInfo(ctx).WorkflowExecution.ID

	savePath := path.Join(cfg.OutputDir, workflowID)

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

func NewActivities(nc *nats.Conn) *activities {
	return &activities{
		nc: nc,
	}
}
