import { Body, Controller, Get, Put, Req, UseGuards } from '@nestjs/common';
import { UserService } from './user.service';
import { Roles } from 'src/common/roles/roles.decorator';
import { RolesGuard } from 'src/common/roles/roles.guard';

@Controller('api/user')
@UseGuards(RolesGuard)
export default class UserController {
  constructor(private readonly userService: UserService) {}
  @Get('me')
  @Roles('USER')
  getProfile(@Req() req) {
    return this.userService.getProfile(req.user.id);
  }

  @Get('orders')
  @Roles('USER')
  getOrders(@Req() req) {
    return this.userService.getOrders(req.user.id);
  }

  @Put('update')
  @Roles('USER')
  updateProfile(@Req() req, @Body() dto: any) {
    return this.userService.updateProfile(req.user.id, req.body);
  }

  @Get('cart')
  @Roles('USER')
  getCart(@Req() req) {
    return this.userService.getCart(req.user.id);
  }
}
