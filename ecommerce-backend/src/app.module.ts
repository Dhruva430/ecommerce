import { Module } from '@nestjs/common';
import { AuthModule } from './auth/auth.module';

import { ProductModule } from './product/product.module';
import { CartModule } from './cart/cart.module';
import { APP_GUARD } from '@nestjs/core';
import { RolesGuard } from './common/roles/roles.guard';

@Module({
  imports: [AuthModule, ProductModule, CartModule],
  providers: [
    {
      provide: APP_GUARD,
      useClass: RolesGuard,
    },
  ],
})
export class AppModule {}
