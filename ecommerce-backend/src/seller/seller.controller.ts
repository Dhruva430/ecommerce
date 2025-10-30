import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  Req,
  UseGuards,
} from '@nestjs/common';
import SellerService from './seller.service';
import { Roles } from 'src/common/roles/roles.decorator';
import { Role } from '@prisma/client';
import { AuthGuard } from 'src/auth/guards/auth.guard';
import { RolesGuard } from 'src/common/roles/roles.guard';

@Controller('api/seller')
@UseGuards(AuthGuard, RolesGuard)
export default class SellerController {
  constructor(private readonly sellerService: SellerService) {}
  @Post('add-products')
  @Roles('SELLER')
  addProduct(@Req() req, @Body() dto: any) {
    return this.sellerService.createProduct(req.user.userId, dto);
  }

  @Get('products')
  @Roles('SELLER')
  getMyProducts(@Req() req) {
    return this.sellerService.getProducts(req.user.userId);
  }

  @Get('orders')
  @Roles('SELLER')
  getSellerOrders(@Req() req) {
    return this.sellerService.getOrders(req.user.userId);
  }

  @Delete('products/:productId')
  @Roles('SELLER')
  deleteProduct(@Param('productId') productId: string) {
    return this.sellerService.deleteProduct(productId);
  }

  @Get('profile/:sellerId')
  async getSellerProfile(@Param('sellerId') sellerId: string) {
    return this.sellerService.getSellerById(sellerId);
  }
  @Put('products/:productId')
  @Roles('SELLER')
  updateProduct(@Req() req, @Body() dto: any) {
    // TODO: Validate dto
    return this.sellerService.updateProduct(
      req.user.userId,
      dto.productId,
      dto,
    );
  }
}
