import { Injectable } from '@nestjs/common';
import { prisma } from '@ecommerce-backend/shared';
@Injectable()
export class AppService {
  async getAllProducts(limit: number = 20, cursor?: string) {
    const products = await prisma.product.findMany({
      take: limit,
      skip: cursor ? 1 : 0,
      cursor: cursor ? { id: cursor } : undefined,
      orderBy: {
        id: 'desc',
      },
    });
    const nextCursor =
      products.length === limit ? products[products.length - 1].id : null;
    return {
      products,
      nextCursor,
    };
  }
  async createProduct() {}
}
