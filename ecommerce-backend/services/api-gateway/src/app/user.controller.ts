import {
  Body,
  Controller,
  Delete,
  Get,
  Inject,
  Param,
  Patch,
  Post,
  Req,
} from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import {
  CreateAddressSchema,
  type CreateAddressDto,
} from '../common/dtos/create-address.dto';
import {
  UpdateAddressSchema,
  UpdateAddressDto,
} from '../common/dtos/update-address.dto';
import { UpdateUserDto } from '../common/dtos/update-user.dto';
import { ZodValidationPipe } from '@ecommerce-backend/shared-dtos';

@Controller('user')
export class UserGatewayController {
  constructor(
    @Inject('USER_SERVICE') private readonly userClient: ClientProxy
  ) {}

  @Get('ping')
  ping() {
    return this.userClient.send('user.ping', {});
  }

  @Get(':id')
  getUser(@Param('id') id: string) {
    return this.userClient.send('user.get', id);
  }

  @Delete(':id')
  deleteUser(@Req() req: Request) {
    return this.userClient.send('user.delete', req.headers);
  }

  @Patch(':id')
  updateUser(@Req() req: Request, @Body() dto: UpdateUserDto) {
    return this.userClient.send('user.update', { headers: req.headers, dto });
  }

  @Get('address')
  getUserAddress(@Req() req: Request) {
    return this.userClient.send('user.address.get', req.headers);
  }

  @Post('address')
  addUserAddress(
    @Req() req: any,
    @Body(new ZodValidationPipe(CreateAddressSchema)) dto: CreateAddressDto
  ) {
    return this.userClient.send('user.address.add', {
      headers: req.headers,
      dto,
    });
  }

  @Patch(':addressId/address')
  updateUserAddress(
    @Param('addressId') addressId: string,
    @Req() req: Request,
    @Body(new ZodValidationPipe(UpdateAddressSchema)) dto: UpdateAddressDto
  ) {
    return this.userClient.send('user.address.update', {
      addressId,
      dto,
      headers: req.headers,
    });
  }

  @Delete(':addressId/address')
  deleteUserAddress(
    @Param('addressId') addressId: string,
    @Req() req: Request
  ) {
    return this.userClient.send('user.address.delete', {
      addressId,
      headers: req.headers,
    });
  }

  @Get('orders')
  getOrders(@Req() req: Request) {
    return this.userClient.send('user.orders', req.headers);
  }

  @Get('reviews')
  getReviews(@Req() req: Request) {
    return this.userClient.send('user.reviews', req.headers);
  }
}
