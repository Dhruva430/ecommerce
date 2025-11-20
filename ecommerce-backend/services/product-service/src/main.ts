/**
 * This is not a production server yet!
 * This is only a minimal backend to get started.
 */

import { Logger } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app/app.module';
import * as dotenv from 'dotenv';
import * as path from 'path';
import helmet from 'helmet';
import compression from 'compression';

async function bootstrap() {
  dotenv.config({
    path: path.join(__dirname, '..', '.env'),
  });
  const app = await NestFactory.create(AppModule);
  app.use(helmet());
  app.use(compression());
  const port = process.env.PORT || 4003;
  await app.listen(port);
  Logger.log(`🚀 Application is running on: http://localhost:${port}/`);
}

bootstrap();
