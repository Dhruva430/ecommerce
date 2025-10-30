import { Module } from '@nestjs/common';
import SellerController from './seller.controller';
import SellerService from './seller.service';
import { JwtModule } from '@nestjs/jwt';
import { RolesGuard } from 'src/common/roles/roles.guard';
import { AuthGuard } from 'src/auth/guards/auth.guard';
import { jwtConstants } from 'libs/constant';

@Module({
  imports: [
    JwtModule.register({
      secret: jwtConstants.secret,
      signOptions: { expiresIn: '1d' },
    }),
  ],
  controllers: [SellerController],
  providers: [SellerService, AuthGuard, RolesGuard],
})
export class SellerModule {}
