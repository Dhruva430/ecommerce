import { Controller, Get, Param } from '@nestjs/common';
import { CartService } from './cart.service';

@Controller('/add-to-cart')
export class CartController {
  constructor(private readonly cartService: CartService) {}

  @Get('/:id')
  addToCart(@Param('id') id: string) {
    return this.cartService.addToCart(id);
  }
}
