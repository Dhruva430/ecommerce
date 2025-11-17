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

  await app.listen(process.env.PORT || 4001);
  console.log(
    'Auth Service running on http://localhost:' + (process.env.PORT || 4001)
  );
}

bootstrap();
