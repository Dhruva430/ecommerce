import { Body, Controller, Get, Inject, Post, Req, Res } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { ZodValidationPipe, Public } from '@ecommerce-backend/shared-dtos';
import { LoginSchema, LoginDto } from '@ecommerce-backend/shared-dtos';
import { SignupSchema, SignupDto } from '@ecommerce-backend/shared-dtos';

@Controller('auth')
export class AuthGatewayController {
  constructor(
    @Inject('AUTH_SERVICE') private readonly authClient: ClientProxy
  ) {}

  @Public()
  @Post('signup')
  signup(@Body(new ZodValidationPipe(SignupSchema)) dto: SignupDto) {
    return this.authClient.send('auth.signup', dto);
  }

  @Public()
  @Get('ping')
  ping() {
    return this.authClient.send('user.ping', {});
  }
  @Public()
  @Post('login')
  login(@Body(new ZodValidationPipe(LoginSchema)) dto: LoginDto) {
    return this.authClient.send('auth.login', dto);
  }

  @Get('me')
  me(@Req() req: Request) {
    return this.authClient.send('auth.me', req.headers);
  }

  @Get('logout')
  logout(@Req() req: Request) {
    return this.authClient.send('auth.logout', req.headers);
  }

  @Get('refresh')
  refresh(@Req() req: Request, @Res({ passthrough: true }) res: Response) {
    return this.authClient.send('auth.refresh', req.headers);
  }
}
