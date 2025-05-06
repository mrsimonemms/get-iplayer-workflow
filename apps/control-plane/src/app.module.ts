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
import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { PrometheusModule } from '@willsoto/nestjs-prometheus';

import config from './config';
import { DownloadsModule } from './downloads/downloads.module';
import { HealthModule } from './health/health.module';
import { MetricsController } from './health/metrics.controller';
import { MessagingModule } from './messaging/messaging.module';
import { TemporalModule } from './temporal/temporal.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: config,
    }),
    PrometheusModule.register({
      // Define the controller so it's not under /v1 namespaces
      controller: MetricsController,
    }),

    DownloadsModule,
    HealthModule,
    TemporalModule,
    MessagingModule,
  ],
})
export class AppModule {}
