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
import { BadRequestException, Inject, Injectable } from '@nestjs/common';
import { Client } from '@temporalio/client';
import { randomBytes } from 'node:crypto';
import { URL } from 'node:url';

import { WORKFLOW_CLIENT } from '../temporal/temporal.providers';

interface IMatcher {
  regex: RegExp;
  extractor?: (url: URL) => string;
}

@Injectable()
export class DownloadsService {
  @Inject(WORKFLOW_CLIENT)
  temporalClient: Client;

  // Take a BBC iPlayer/Sounds URL and download it as a Temporal workflow
  async downloadFromURL(inputURL: string): Promise<string> {
    // Extract the programme ID from the URL
    const pid = this.parseURLToPID(inputURL);

    // Generate a random ID
    const id = randomBytes(3).toString('hex');

    const workflow = await this.temporalClient.workflow.start(
      'DownloadBBCProgramme',
      {
        taskQueue: 'downloadByPID',
        workflowId: `download-pid-${pid}-${id}`,
        args: [
          {
            programmeID: pid,
          },
        ],
      },
    );

    return workflow.workflowId;
  }

  // Extract the BBC PID from a URL
  // @link https://en.wikipedia.org/wiki/BBC_Programme_Identifier
  parseURLToPID(inputURL: string): string {
    const pidRegex = '[a-z0-9]{7,}';
    const matchers: IMatcher[] = [
      // Order is important
      {
        regex: new RegExp(`^${pidRegex}-[\\w]+-\\d-(${pidRegex})`),
        extractor(url: URL): string {
          return url.searchParams.get('seriesId') ?? '';
        },
      },
      {
        regex: new RegExp(`^/sounds/play/(${pidRegex})`),
      },
      {
        regex: new RegExp(`^/programmes/(${pidRegex})`),
      },
      {
        regex: new RegExp(`^/iplayer/episode/(${pidRegex})`),
      },
    ];

    const url = URL.parse(inputURL);
    if (!url) {
      throw new BadRequestException('Cannot parse as URL');
    }

    const matcher = matchers.find((m) => {
      const s = m.extractor ? m.extractor(url) : url.pathname;
      return m.regex.test(s);
    });

    const inputString = matcher?.extractor
      ? matcher.extractor(url)
      : url.pathname;

    const matches = inputString.match(matcher?.regex ?? '');

    if (!matches || !matches[1]) {
      throw new Error('Cannot extract the programme ID');
    }

    return matches[1];
  }
}
