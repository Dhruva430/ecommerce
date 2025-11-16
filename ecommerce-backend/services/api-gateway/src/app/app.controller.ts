import { Controller, Get, Req } from '@nestjs/common';

@Controller()
export class AppController {
  @Get('health')
  getHealth() {
    return { status: 'ok' };
  }

  @Get('me')
  getProfile(@Req() req: Request) {
    return {
      message: 'Protected route accessed successfully',
      user: (req as any).user,
    };
  }
}
