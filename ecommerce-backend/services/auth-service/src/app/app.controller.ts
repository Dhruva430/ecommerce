import { Controller } from '@nestjs/common';
import { MessagePattern } from '@nestjs/microservices';

@Controller()
export class AppController {
  @MessagePattern('auth.login')
  login() {
    return { message: 'Logged in' };
  }
}
