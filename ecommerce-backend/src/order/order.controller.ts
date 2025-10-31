import { Get, Post, Query, Req } from '@nestjs/common';
import { OrderService } from './order.service';

export default class OrderController {
  constructor(private readonly orderService: OrderService) {}

  @Post('buy-order')
  buyOrder(@Req() req, variantId: string, quantity: number, address: string) {
    return this.orderService.buyOrder(
      req.user.id,
      variantId,
      quantity,
      address,
    );
  }
}
