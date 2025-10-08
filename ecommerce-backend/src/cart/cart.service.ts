import { Injectable } from '@nestjs/common';
import productData from '../data/products.json';
import { prisma } from 'libs/primsa';
@Injectable()
export class CartService {
  private products = productData;

  async addToCart(id: string) {
    const product = this.products.find((product) => product.id === id);
    if (!product) {
      return { message: 'Product not found' };
    }
    if (!product.price || !product.price.payable) {
      return { message: 'Product price not available' };
    }
    const parsedAmount = parseFloat(product.price.payable);
    await prisma.cart.upsert({
      where: {
        userId_productId: {
          userId: '66f55c63b5e0a5c6e3dbfabc',
          productId: id,
        },
      },
      update: {
        amount: { increment: parsedAmount },
      },
      create: {
        userId: '66f55c63b5e0a5c6e3dbfabc',
        productId: id,
        amount: parsedAmount,
      },
    });

    return { message: 'Product added to cart successfully' };
  }
}
