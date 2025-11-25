import { Controller } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { AppService } from './app.service';
import { SignupDto } from '@ecommerce-backend/shared-dtos';
import { LoginDto } from '@ecommerce-backend/shared-dtos';

@Controller()
export class AuthController {
  constructor(private readonly appService: AppService) {}

  @MessagePattern('auth.signup')
  signup(@Payload() dto: SignupDto) {
    return this.appService.signup(dto);
  }

  @MessagePattern('auth.login')
  login(@Payload() dto: LoginDto) {
    return this.appService.login(dto);
  }

  @MessagePattern('auth.me')
  me(@Payload() headers: Record<string, string>) {
    return this.appService.me(headers);
  }

  @MessagePattern('auth.logout')
  logout(@Payload() headers: Record<string, string>) {
    return this.appService.logout(headers);
  }

  @MessagePattern('auth.refresh')
  refresh(@Payload() headers: Record<string, string>, res: Response) {
    return this.appService.refresh(headers, res);
  }
}
