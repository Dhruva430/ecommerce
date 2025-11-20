import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  Req,
} from '@nestjs/common';
import { AppService } from './app.service';
import {
  createAddressSchema,
  type createAddressDto,
} from '../common/dtos/create-address.dto';
import {
  updateAddressSchema,
  type updateAddressDto,
} from '../common/dtos/update-address.dto';
import { type updateUserDto } from '../common/dtos/update-user.dto';
import { ZodValidationPipe } from '@ecommerce-backend/shared';

@Controller('user')
export class AppController {
  constructor(private readonly userService: AppService) {}
  // ------------------- User Management -------------------
  @Get('ping')
  ping() {
    return 'pong';
  }
  @Get(':id')
  getUser(@Param('id') id: string) {
    return this.userService.getUser(id);
  }
  @Delete(':id')
  deleteUser(@Req() req: Request) {
    return this.userService.deleteUser(req);
  }
  @Patch(':id')
  updateUser(@Req() req: Request, dto: updateUserDto) {
    return this.userService.updateUser(req, dto);
  }
  // ------------------ Address Management ------------------
  @Get('address')
  getUserAddress(@Req() req: Request) {
    return this.userService.getUserAddress(req);
  }
  @Post('address')
  addUserAddress(
    @Req() req: Request,
    @Body(new ZodValidationPipe(createAddressSchema)) dto: createAddressDto
  ) {
    return this.userService.addUserAddress(req, dto);
  }
  @Patch(':addressId/address')
  updateUserAddress(
    @Param('addressId') addressId: string,
    @Req() req: Request,
    @Body(new ZodValidationPipe(updateAddressSchema)) dto: updateAddressDto
  ) {
    return this.userService.updateUserAddress(addressId, dto, req);
  }
  @Delete(':addressId/address')
  deleteUserAddress(
    @Param('addressId') addressId: string,
    @Req() req: Request
  ) {
    return this.userService.deleteUserAddress(req, addressId);
  }

  // ------------------- Wishlist Management -------------------
  // @Get(':id/wishlist')
  // getUserWishlist(@Param('id') id: string) {
  //   return this.userService.getUserWishlist(id);
  // }
  // @Post(':id/wishlist')
  // addUserWishlistItem(@Param('id') id: string) {
  //   return this.userService.addUserWishlistItem(id);
  // }
  // @Delete(':id/wishlist')
  // deleteUserWishlistItem(@Param('id') id: string) {
  //   return this.userService.deleteUserWishlistItem(id);
  // }
  // ------------------- User Activity -------------------
  @Get('orders')
  getUserOrders(@Req() req: Request) {
    return this.userService.getUserOrders(req);
  }
  @Get('reviews')
  getUserReviews(@Req() req: Request) {
    return this.userService.getUserReviews(req);
  }
}
