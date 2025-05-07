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

	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/log"
)

type BBCProgrammeAPI struct {
	Programme struct {
		Position     int    `json:"position"`
		Title        string `json:"title"`
		DisplayTitle struct {
			Title string `json:"title"`
		} `json:"display_title"`
		Parent struct {
			Programme struct {
				Position int    `json:"position"`
				Title    string `json:"title"`
			} `json:"programme"`
		} `json:"parent"`
	} `json:"programme"`
}

func (p *BBCProgrammeAPI) GetFileName(ext string) string {
	episodeTitle := removeNonAlnum(p.Programme.Title)
	showTitle := removeNonAlnum(p.Programme.DisplayTitle.Title)
	episodeNumber := p.Programme.Position
	seriesNumber := p.Programme.Parent.Programme.Position

	var name string
	if episodeNumber == 0 && seriesNumber == 0 {
		// Treat as a single programme
		name = episodeTitle
	} else {
		// Treat as part of a series
		name = fmt.Sprintf(
			"%s - s%se%s - %s",
			showTitle,
			leftPad(seriesNumber),
			leftPad(episodeNumber),
			episodeTitle,
		)
	}

	name += ext

	return name
}

type Download struct {
	ProgrammeID string
}

type DownloadBBCProgrammeResult struct{}

type DownloadByPIDResult struct {
	ProgrammeID string
	SavePath    string
	Files       []string
}

type DownloadedProgramme struct {
	ProgrammeID string
	SavePath    string
	File        string
	TargetName  string
}

type ParseDownloadedProgrammeResult struct {
	ProgrammeID string
	Bucket      *minio.UploadInfo
}

type ProgrammeNameResult struct {
	Name string
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
