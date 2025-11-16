import {
  Module,
  MiddlewareConsumer,
  NestModule,
  RequestMethod,
} from '@nestjs/common';
import { RateLimitMiddleware } from '@ecommerce-backend/shared';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AuthMiddleware } from '../common/middleware/auth.middleware';
import { ProxyService } from '../common/proxy/proxy.service';

@Module({
  imports: [],
  controllers: [AppController],
  providers: [AppService, ProxyService],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(RateLimitMiddleware).forRoutes('*');
    consumer
      .apply(AuthMiddleware)
      .exclude(
        { path: 'health', method: RequestMethod.ALL },
        { path: 'auth/login', method: RequestMethod.ALL },
        { path: 'auth/signup', method: RequestMethod.ALL },
        { path: 'auth/request-otp', method: RequestMethod.ALL }
      )
      .forRoutes('*');
  }
}
