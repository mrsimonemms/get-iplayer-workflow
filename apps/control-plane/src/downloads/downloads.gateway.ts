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
import { Inject, Logger } from '@nestjs/common';
import {
  ConnectedSocket,
  MessageBody,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { DefaultEventsMap, Server, Socket } from 'socket.io';

import { MessagingService } from '../messaging/messaging.service';

export type SocketClient = Socket<
  DefaultEventsMap,
  DefaultEventsMap,
  DefaultEventsMap,
  SocketData
>;

export interface SocketData {
  events: Map<string, any>;
}

@WebSocketGateway({ transports: ['websocket'] })
export class DownloadsGateway {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(MessagingService)
  messagingClient: MessagingService;

  @WebSocketServer()
  server: Server<
    DefaultEventsMap,
    DefaultEventsMap,
    DefaultEventsMap,
    SocketData
  >;

  @SubscribeMessage('download:msg')
  async handleDownload(
    @ConnectedSocket()
    client: Socket<
      DefaultEventsMap,
      DefaultEventsMap,
      DefaultEventsMap,
      SocketData
    >,
    @MessageBody() workflowId: string,
  ) {
    const sub = this.messagingClient.subscribe(`log.${workflowId}`);

    for await (const m of sub) {
      client.emit('log', m.string());
    }
  }
}
