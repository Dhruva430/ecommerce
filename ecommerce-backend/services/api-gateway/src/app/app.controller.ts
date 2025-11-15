import { Public } from '@ecommerce-backend/shared';
import { Controller, Get } from '@nestjs/common';

@Controller()
export class AppController {
  @Get('health')
  @Public()
  getHealth() {
    return { status: 'ok' };
  }
}
