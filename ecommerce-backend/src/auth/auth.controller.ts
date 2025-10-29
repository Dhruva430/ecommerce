import { Body, Controller, Get, Post, Req } from '@nestjs/common';
import { type CreateUserDto } from 'src/common/dto/create-user.dto';
import AuthService from './auth.service';

@Controller('api/auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post('register')
  register(@Body() dto: CreateUserDto) {
    return this.authService.register(dto);
  }

  @Post('login')
  login(@Body() dto: CreateUserDto) {
    return this.authService.login(dto);
  }
}
