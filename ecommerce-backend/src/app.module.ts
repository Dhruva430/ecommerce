import { Module } from '@nestjs/common';
import { AuthModule } from './auth/auth.module';

import { ProductModule } from './product/product.module';
import { CartModule } from './cart/cart.module';
import { APP_GUARD } from '@nestjs/core';
import { RolesGuard } from './common/roles/roles.guard';
import { UserModule } from './user/user.module';
import { SellerModule } from './seller/seller.module';

@Module({
  imports: [AuthModule, ProductModule, CartModule, UserModule, SellerModule],
  providers: [
    {
      provide: APP_GUARD,
      useClass: RolesGuard,
    },
  ],
})
export class AppModule {}
