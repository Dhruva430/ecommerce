import { Injectable } from '@nestjs/common';

import products from '../data/products.json';
import { prisma } from 'libs/primsa';
import { Category } from '@prisma/client';
@Injectable()
export class ProductService {
  async getProducts(limit: number = 10, offset: number = 0) {
    const product = await prisma.product.findMany({
      take: limit,
      skip: offset,
    });

    return product;
  }
  async restoreProducts() {
    for (const product of products) {
      await prisma.product.create({
        data: {
          id: product.id,
          title: product.title,
          description: product.description,
          category: Category.ELECTRONICS,
          imageUrl: product.imageUrl,
          price: product.price,
          discounted: product.discounted,
          sellerId: 'd0090758-0148-44df-b7db-582b9dd5f556',
          productVariant: {
            create: product.ProductVariant.map((variant) => ({
              id: variant.id,
              size: variant.size,
              description: variant.description,
              title: variant.title,
              price: variant.price,
              stock: variant.stock,

              variantAttribute: {
                create: variant.VariantAttribute.map((attr) => ({
                  id: attr.id,
                  name: attr.name,
                  value: attr.value,
                })),
              },

              variantImages: {
                create: variant.VariantImages.map((img) => ({
                  id: img.id,
                  imageUrl: img.imageUrl,
                  position: img.position,
                })),
              },
            })),
          },
        },
      });
    }
    return { message: 'Products inserted successfully' };
  }
  async getProductById(id: string) {
    const product = await prisma.product.findMany({
      where: { id },
    });
    if (!product) {
      throw new Error('Product not found');
    }
    return product;
  }
  async getProductByCategory(category: string) {
    const product = await prisma.product.findMany({
      where: { category: Category.ELECTRONICS },
    });
    if (!product) {
      throw new Error('Product not found');
    }
    return product;
  }
}
