import { Body, Controller, Post } from '@nestjs/common';
import { AppService } from './app.service';
import { ZodValidationPipe } from '@ecommerce-backend/shared';
import { LoginSchema, type LoginDto } from '../common/dtos/login';
import { SignupSchema, type SignupDto } from '../common/dtos/signup';
import {
  RequestOtpSchema,
  type RequestOtpDto,
} from '../common/dtos/request-otp';

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

  @Post('request-otp')
  requestOtp(
    @Body(new ZodValidationPipe(RequestOtpSchema)) dto: RequestOtpDto
  ) {
    return this.authService.requestOtp(dto);
  }

  @Post('verify-otp')
  verifyOtp(@Body('email') email: string, @Body('otp') otp: string) {
    return this.authService.verifyOtp(email, otp);
  }
}
