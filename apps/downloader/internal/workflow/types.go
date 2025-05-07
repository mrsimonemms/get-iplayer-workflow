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

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/log"
)

type Download struct {
	ProgrammeID string
}

type DownloadByPIDResult struct {
	ProgrammeID string
	SavePath    string
	Files       []string
}

type DownloadedProgramme struct {
	ProgrammeID string
	SavePath    string
	File        string
}

type UploadedProgramme struct{}

type streamOutput struct {
	logger     log.Logger
	nc         *nats.Conn
	workflowID string
}

func (s *streamOutput) Write(p []byte) (n int, err error) {
	s.logger.Debug("New data received", "msg", string(p), "workflowID", s.workflowID)
	if err := s.nc.Publish(fmt.Sprintf("log.%s", s.workflowID), p); err != nil {
		s.logger.Error("Error emitting message to NATS", "error", err)
		return 0, err
	}
	return len(p), nil
}
