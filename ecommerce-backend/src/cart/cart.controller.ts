import {
  Controller,
  Delete,
  Get,
  Param,
  ParseIntPipe,
  Query,
  UseGuards,
} from '@nestjs/common';
import { CartService } from './cart.service';
import { AuthGuard } from 'src/auth/guards/auth.guard';
import { RolesGuard } from 'src/common/roles/roles.guard';
import { Roles } from 'src/common/roles/roles.decorator';

@Controller('/cart')
@UseGuards(AuthGuard, RolesGuard)
export class CartController {
  constructor(private readonly cartService: CartService) {}

  @Get('add')
  @Roles('USER')
  addToCart(
    @Query('productId') productId: string,
    @Query('amount', ParseIntPipe) amount: number,
  ) {
    return this.cartService.addToCart(productId, amount);
  }

  @Delete('/remove')
  @Roles('USER')
  removeFromCart(@Query('productId') productId: string) {
    return this.cartService.removeFromCart(productId);
  }
}
