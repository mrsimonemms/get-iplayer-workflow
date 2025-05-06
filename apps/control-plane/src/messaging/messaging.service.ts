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
import {
  NatsConnection,
  Subscription,
  SubscriptionOptions,
} from '@nats-io/nats-core';
import {
  Inject,
  Injectable,
  Logger,
  OnApplicationShutdown,
} from '@nestjs/common';

import { CONNECTION } from './messaging.providers';

@Injectable()
export class MessagingService implements OnApplicationShutdown {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(CONNECTION)
  private connection: NatsConnection;

  subscribe(subject: string, opts?: SubscriptionOptions): Subscription {
    return this.connection.subscribe(subject, opts);
  }

  async onApplicationShutdown(): Promise<void> {
    this.logger.log('Disconnecting from NATS server');
    const done = this.connection.closed();
    await this.connection.close();
    const err = await done;
    if (err) {
      this.logger.error('Error closing NATS connection', err);
    } else {
      this.logger.debug('Successfully closed NATS server');
    }
  }
}
