import { Module } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { AuthModule } from './auth/auth.module';
import { ProductModule } from './product/product.module';
import { CartModule } from './cart/cart.module';
import { UserModule } from './user/user.module';
import { SellerModule } from './seller/seller.module';
import { RolesGuard } from './common/roles/roles.guard';
import { AuthGuard } from './auth/guards/auth.guard';

@Module({
  imports: [AuthModule, ProductModule, CartModule, UserModule, SellerModule],
})
export class AppModule {}
