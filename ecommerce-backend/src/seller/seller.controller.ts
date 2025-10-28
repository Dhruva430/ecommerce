import { Body, Controller, Delete, Get, Post, Put, Req } from '@nestjs/common';
import SellerService from './seller.service';
import { Roles } from 'src/common/roles/roles.decorator';
import { Role } from '@prisma/client';

@Controller('api/seller')
export default class SellerController {
  constructor(private readonly sellerService: SellerService) {}
  @Post('products')
  @Roles(Role.SELLER)
  addProduct(@Req() req, @Body() dto: any) {
    return this.sellerService.createProduct(req.user.userId, dto);
  }

  @Get('products')
  @Roles(Role.SELLER)
  getMyProducts(@Req() req) {
    return this.sellerService.getProducts(req.user.userId);
  }

  @Get('orders')
  @Roles(Role.SELLER)
  getSellerOrders(@Req() req) {
    return this.sellerService.getOrders(req.user.userId);
  }

  @Delete('products/:productId')
  @Roles(Role.SELLER)
  deleteProduct(@Req() req, @Body() dto: any) {
    return this.sellerService.deleteProduct(req.user.userId, dto.productId);
  }

  @Get('profile/:sellerId')
  async getSellerProfile(@Req() req) {
    const sellerId = req.params.sellerId;
    return this.sellerService.getSellerById(sellerId);
  }
  @Put('products/:productId')
  @Roles(Role.SELLER)
  updateProduct(@Req() req, @Body() dto: any) {
    return this.sellerService.updateProduct(
      req.user.userId,
      dto.productId,
      dto,
    );
  }
}
