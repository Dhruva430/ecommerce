import { Controller, Get, Post, Body, Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { GatewayService } from './gateway.service';
import { type CreateUserDto } from '@shared/dto/create-user.dto';

@Controller('gateway')
export class GatewayController {
  constructor(
    private readonly gatewayService: GatewayService,
    @Inject('USER_SERVICE') private readonly userClient: ClientProxy,
  ) {}

  @Get('health')
  getHealth() {
    return { status: 'API Gateway running smoothly ✅' };
  }

  @Post('user/register')
  async registerUser(@Body() dto: CreateUserDto) {
    return this.userClient.send({ cmd: 'register_user' }, dto);
  }
}
