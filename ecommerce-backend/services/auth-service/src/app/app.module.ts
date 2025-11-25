import { Module } from '@nestjs/common';
import { AuthController } from './app.controller';
import { AppService } from './app.service';
import { JwtModule } from '@nestjs/jwt';
import { TokenUtil } from '../common/utils/token';
import { config } from 'shared/src/lib/config';

@Module({
  imports: [
    JwtModule.register({
      secret: config.JWT_SECRET,
      signOptions: { expiresIn: '15m' }, // default for access token
    }),
  ],
  controllers: [AuthController],
  providers: [AppService, TokenUtil],
  exports: [AppService, TokenUtil],
})
export class AppModule {}
