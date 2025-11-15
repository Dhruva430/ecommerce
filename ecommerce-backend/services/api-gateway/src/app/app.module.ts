import { Module, MiddlewareConsumer, NestModule } from '@nestjs/common';
import { RateLimitMiddleware, ProxyService } from '@ecommerce-backend/shared';
import { AppController } from './app.controller';
import { AppService } from './app.service';

@Module({
  imports: [],
  controllers: [AppController],
  providers: [AppService, ProxyService],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(RateLimitMiddleware).forRoutes('*');
  }
}
