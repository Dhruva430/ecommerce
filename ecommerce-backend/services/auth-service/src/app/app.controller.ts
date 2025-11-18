import { Body, Controller, Get, Post, Req } from '@nestjs/common';
import { AppService } from './app.service';
import { ZodValidationPipe } from '@ecommerce-backend/shared';
import { LoginSchema, type LoginDto } from '../common/dtos/login';

import { type SignupDto, SignupSchema } from '../common/dtos/signup';

@Controller('auth')
export class AppController {
  constructor(private readonly authService: AppService) {}

  @Post('signup')
  signup(@Body(new ZodValidationPipe(SignupSchema)) dto: SignupDto) {
    return this.authService.signup(dto);
  }

  @Post('login')
  login(@Body(new ZodValidationPipe(LoginSchema)) dto: LoginDto) {
    return this.authService.login(dto);
  }
  @Get('me')
  me(@Req() req: Request) {
    return this.authService.me(req);
  }
  @Get('logout')
  logout(@Req() req: Request) {
    return this.authService.logout(req);
  }
}
