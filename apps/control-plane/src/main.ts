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
  ClassSerializerInterceptor,
  ConsoleLogger,
  Logger,
  ValidationPipe,
  VersioningType,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { NestFactory, Reflector } from '@nestjs/core';
import {
  FastifyAdapter,
  NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

import { AppModule } from './app.module';
import loggerConfig from './config/logger';
import { ISocketConfig } from './config/socket';
import { WsExceptionFilter } from './lib/customFilters';
import { RedisIoAdapter } from './lib/redisAdapter';

const logger = new Logger('bootstrap');

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
    { logger: new ConsoleLogger(loggerConfig()) },
  );

  const config = app.get(ConfigService);

  app
    .enableShutdownHooks()
    .enableVersioning({
      type: VersioningType.URI,
      defaultVersion: '1',
    })
    .useGlobalFilters(new WsExceptionFilter())
    .useGlobalPipes(
      new ValidationPipe({
        transform: true,
        whitelist: true,
      }),
    )
    .useGlobalInterceptors(new ClassSerializerInterceptor(app.get(Reflector)));

  // Add Swagger documentation
  const docBuilderConfig = new DocumentBuilder()
    .setTitle(process.env.npm_package_name!)
    .setDescription(
      'Temporal workflow to download and sort stuff from the iPlayer',
    )
    .setVersion(process.env.name_package_version ?? 'dev')
    .addBearerAuth()
    .build();

  const documentFactory = SwaggerModule.createDocument(app, docBuilderConfig);

  SwaggerModule.setup('api', app, documentFactory);

  // Configure Redis for SocketIO
  const socketRedis = config.getOrThrow<ISocketConfig>('socket');
  if (socketRedis.redis.enabled) {
    logger.debug('Connect socket redis adapter');
    const redisIoAdapter = new RedisIoAdapter(app);
    await redisIoAdapter.connectToRedis(socketRedis.redis.opts);
    app.useWebSocketAdapter(redisIoAdapter);
  }

  await app.listen(
    config.getOrThrow<number>('server.port'),
    config.getOrThrow('server.host'),
  );
}

bootstrap().catch((err: Error) => {
  /* Unlikely to get to here but a final catchall */
  logger.fatal(err.stack);
});
