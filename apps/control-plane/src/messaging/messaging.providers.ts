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
import { ConnectionOptions } from '@nats-io/nats-core';
import { connect } from '@nats-io/transport-node';
import { Provider } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

export const CONNECTION = Symbol('CONNECTION');

export const messagingProviders: Provider[] = [
  {
    inject: [ConfigService],
    provide: CONNECTION,
    useFactory: (cfg: ConfigService) => {
      const opts: ConnectionOptions = {
        servers: cfg.getOrThrow('nats.url'),
      };
      return connect(opts);
    },
  },
];
