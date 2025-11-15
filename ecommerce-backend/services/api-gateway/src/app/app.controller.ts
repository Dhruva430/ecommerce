import { AuthGuard } from '@ecommerce-backend/shared';
import { Controller, Get, Req, UseGuards } from '@nestjs/common';

@Controller()
export class AppController {
  @Get('health')
  getHealth() {
    return { status: 'ok' };
  }

  @Get('me')
  @UseGuards(AuthGuard)
  getProfile(@Req() req: Request) {
    return {
      message: 'Protected route accessed successfully',
      user: (req as any).user,
    };
  }
}
