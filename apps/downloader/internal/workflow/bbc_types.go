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

import "fmt"

type BBCProgrammeAPIProgrammeMediaType string

const BBCProgrammeAPIProgrammeMediaTypeAudio BBCProgrammeAPIProgrammeMediaType = "audio"

type BBCProgrammeAPIDisplayTitle struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type BBCProgrammeAPIParent struct {
	Programme BBCProgrammeAPIParentProgramme `json:"programme"`
}

type BBCProgrammeAPIParentProgramme struct {
	Position int    `json:"position"`
	Title    string `json:"title"`
}

type BBCProgrammeAPIProgramme struct {
	Position       int                               `json:"position"`
	Title          string                            `json:"title"`
	DisplayTitle   BBCProgrammeAPIDisplayTitle       `json:"display_title"`
	Parent         BBCProgrammeAPIParent             `json:"parent"`
	MediaType      BBCProgrammeAPIProgrammeMediaType `json:"media_type"`
	ShortSynopsis  string                            `json:"short_synopsis"`
	MediumSynopsis string                            `json:"medium_synopsis"`
	LongSynopsis   string                            `json:"long_synopsis"`
}

type BBCProgrammeAPI struct {
	Programme BBCProgrammeAPIProgramme `json:"programme"`
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
