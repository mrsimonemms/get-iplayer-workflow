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
import { Controller, Get, Inject } from '@nestjs/common';

import { DownloadsService } from './downloads.service';

@Controller('download')
export class DownloadsController {
  @Inject(DownloadsService)
  private readonly downloadsService: DownloadsService;

  @Get()
  download() {
    return this.downloadsService.downloadFromURL(
      'https://www.bbc.co.uk/sounds/play/m0008bbc',
    );
  }
}
