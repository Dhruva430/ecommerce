import { NestFactory } from '@nestjs/core';
import { AppModule } from './app/app.module';
import helmet from 'helmet';
import compression from 'compression';
import { ProxyService } from './common/proxy/proxy.service';
import * as dotenv from 'dotenv';
import * as path from 'path';

async function bootstrap() {
  dotenv.config({
    path: path.join(__dirname, '..', '.env'),
  });

  const app = await NestFactory.create(AppModule);
  app.setGlobalPrefix('api');
  app.use(helmet());
  app.use(compression());

  const proxyService = app.get(ProxyService);

  const expressApp = app.getHttpAdapter().getInstance();
  proxyService.register(expressApp);

  await app.listen(process.env.PORT || 3000);
  console.log('API Gateway running on http://localhost:3000');
}
bootstrap();
