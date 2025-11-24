import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { AuthGatewayController } from './auth.controller';
import { AuthGuard, RolesGuard } from '@ecommerce-backend/shared-dtos';
import { APP_GUARD } from '@nestjs/core';
import { UserGatewayController } from './user.controller';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'AUTH_SERVICE',
        transport: Transport.REDIS,
        options: { host: 'localhost', port: 6379 },
      },
      {
        name: 'USER_SERVICE',
        transport: Transport.REDIS,
        options: { host: 'localhost', port: 6379 },
      },
    ]),
  ],
  controllers: [AuthGatewayController, UserGatewayController],
  providers: [
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
