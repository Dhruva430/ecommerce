import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { APP_GUARD } from '@nestjs/core';
import { AuthGuard, RolesGuard } from '@ecommerce-backend/shared';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { AppModule as UserModule } from '../../../user-service/src/app/app.module';
// import { ProductModule } from '../modules/product/product.module';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'USER_SERVICE',
        transport: Transport.TCP,
        options: {
          host: process.env.USER_SERVICE_HOST || 'user-service',
          port: 4001,
        },
      },
      // {
      //   name: 'PRODUCT_SERVICE',
      //   transport: Transport.TCP,
      //   options: {
      //     host: process.env.PRODUCT_SERVICE_HOST || 'product-service',
      //     port: 4002,
      //   },
      // },
    ]),
    UserModule,
    // ProductModule,
  ],
  controllers: [AppController],
  providers: [
    AppService,
    {
      provide: APP_GUARD,
      useClass: AuthGuard,
    },
    {
      provide: APP_GUARD,
      useClass: RolesGuard,
    },
  ],
})
export class AppModule {}
