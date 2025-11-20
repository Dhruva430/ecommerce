import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { TokenUtil } from '../common/utils/token';
import { JwtModule } from '@nestjs/jwt';

@Module({
  imports: [
    AppModule,
    JwtModule.register({
      secret: process.env.JWT_SECRET,
      signOptions: { expiresIn: '15m' },
    }),
  ],
  controllers: [AppController],
  providers: [AppService, TokenUtil],
})
export class AppModule {}
