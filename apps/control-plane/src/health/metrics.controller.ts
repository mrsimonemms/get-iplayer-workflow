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
import { Controller, Get, Res, VERSION_NEUTRAL } from '@nestjs/common';
import { FastifyReply } from 'fastify';
import * as client from 'prom-client';

@Controller({
  version: VERSION_NEUTRAL,
})
export class MetricsController {
  @Get()
  index(@Res({ passthrough: true }) response: FastifyReply): Promise<string> {
    response.header('Content-Type', client.register.contentType);
    return client.register.metrics();
  }
}
