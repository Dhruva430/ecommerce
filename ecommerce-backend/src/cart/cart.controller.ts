import { Controller, Get, Param } from '@nestjs/common';
import { CartService } from './cart.service';

@Controller('/add-to-cart')
export class CartController {
  constructor(private readonly cartService: CartService) {}

  @Get('/:productId')
  addToCart(@Param('productId') productId: string) {
    return this.cartService.addToCart(productId);
  }
}
