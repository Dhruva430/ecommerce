import { Controller, Param } from '@nestjs/common';
import { CartService } from './cart.service';

@Controller('/cart')
export class CartController {
  constructor(private readonly cartService: CartService) {}

  AddToCart(@Param('/:id') id: string) {
    return this.cartService.AddToCart(id);
  }
}
