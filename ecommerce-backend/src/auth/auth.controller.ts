import { Controller, Get } from '@nestjs/common';

@Controller('api/auth')
export class AuthController {
  @Get()
  getStatus(): string {
    return 'Auth service is running';
  }
}
