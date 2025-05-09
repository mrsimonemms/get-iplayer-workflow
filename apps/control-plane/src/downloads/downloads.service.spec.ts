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
import { BadRequestException } from '@nestjs/common';
import { Test, TestingModule } from '@nestjs/testing';

import { WORKFLOW_CLIENT } from '../temporal/temporal.providers';
import { DownloadsService } from './downloads.service';

describe('DownloadsService', () => {
  let service: DownloadsService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [DownloadsService],
    })
      .useMocker((token) => {
        if (token === WORKFLOW_CLIENT) {
          return {
            workflow: {
              start: jest.fn(),
            },
          };
        }
      })
      .compile();

    service = module.get<DownloadsService>(DownloadsService);

    expect(service).toBeDefined();
  });

  describe('#parseURLToID', () => {
    it('should parse the iPlayer and Sounds URLs to IDs', () => {
      const data: { url: string; id: string }[] = [
        {
          url: 'https://www.bbc.co.uk/sounds/play/m0008bbc',
          id: 'm0008bbc',
        },
        {
          url: 'https://www.bbc.co.uk/programmes/b09fb6n5/episodes/player',
          id: 'b09fb6n5',
        },
        {
          url: 'https://www.bbc.co.uk/iplayer/episode/m002b1ls/beyond-paradise-series-3-episode-4',
          id: 'm002b1ls',
        },
        {
          url: 'https://www.bbc.co.uk/iplayer/episode/m002b1ls/beyond-paradise-series-3-episode-4?seriesId=m001jg5h-structural-2-m001xgqh',
          id: 'm001xgqh',
        },
        {
          url: 'https://www.bbc.co.uk/iplayer/episode/b0bc5spx/the-kings-speech',
          id: 'b0bc5spx',
        },
        {
          url: 'https://www.bbc.co.uk/iplayer/episodes/b006t0qx/new-tricks',
          id: 'b006t0qx',
        },
      ];

      data.forEach((item) => {
        expect(service.parseURLToPID(item.url)).toEqual(item.id);
      });
    });

    it('should throw an error if an invalid URL sent', () => {
      expect(() => service.parseURLToPID('')).toThrow(
        new BadRequestException('Cannot parse as URL'),
      );
    });

    it('should throw an error if not matcher found', () => {
      expect(() => service.parseURLToPID('https://www.simonemms.com')).toThrow(
        new BadRequestException('Cannot extract the programme ID'),
      );
    });
  });
});
