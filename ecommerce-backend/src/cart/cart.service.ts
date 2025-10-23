import { Injectable } from '@nestjs/common';
import productData from '../data/products.json';
import { prisma } from 'libs/primsa';
@Injectable()
export class CartService {
  async addToCart(productId: string) {
    const cartItem = await prisma.cart.create({
      data: {
        amount: 1,
        productId: productId,
        userId: 'user-123',
      },
    });
    if (!cartItem) {
      throw new Error('Failed to add product to cart');
    }
    return { message: `${cartItem.productId} Added in your Cart. ` };
  }
  async removeFromCart(productId: string) {
    const deletedItem = await prisma.cart.deleteMany({
      where: { productId },
    });
    if (deletedItem.count === 0) {
      throw new Error('Product not found in cart');
    }
    return { message: `Product ${productId} removed from cart.` };
  }
}
