import { Get } from '@nestjs/common';

export class AuthController {
  @Get()
  getHello(): string {
    return 'Hello World!';
  }
}
