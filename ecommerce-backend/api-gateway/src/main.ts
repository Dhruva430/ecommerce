import { NestFactory } from '@nestjs/core';
import { GatewayModule } from './gateway.module';
import { ZodValidationPipe } from '../../shared/pipes/zod-validation.pipe';

async function bootstrap() {
  const app = await NestFactory.create(GatewayModule);

  // Global Zod Validation
  app.useGlobalPipes(new ZodValidationPipe());

  const port = process.env.PORT || 3000;
  await app.listen(port);
  console.log(`🚀 Service running on port ${port}`);
}
bootstrap();
