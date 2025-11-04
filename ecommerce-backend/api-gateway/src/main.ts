import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import { ZodValidationPipe } from '../../shared/pipes/zod-validation.pipe';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // Global Zod Validation
  app.useGlobalPipes(new ZodValidationPipe());

  const port = process.env.PORT || 3000;
  await app.listen(port);
  console.log(`🚀 Service running on port ${port}`);
}
bootstrap();
