import { Injectable } from '@nestjs/common';

@Injectable()
export class CartService {
  AddToCart(id: string): string {
    return `This action adds product with id: ${id} to the cart`;
  }
}
