import { Injectable } from '@nestjs/common';
import { prisma } from 'lib/primsa';
import { type CreateProductDto } from './dto/product.dto';

@Injectable()
export class ProductService {
  async getProduct(id: string) {
    const parsedId = parseInt(id, 10);
    if (isNaN(parsedId)) {
      throw new Error('Invalid product ID');
    }
    const product = await prisma.product.findUnique({
      where: { id: parsedId },
    });
    return product;
  }
  async createProduct(createProductDto: CreateProductDto) {
    const product = await prisma.product.create({
      data: createProductDto,
    });
    return {
      message: 'Product created successfully',
      product: product,
    };
  }
  async deleteProduct(id: string) {
    const parsedInt = parseInt(id, 10);
    await prisma.product.delete({
      where: { id: parsedInt },
    });
    return {
      message: 'Product deleted successfully',
    };
  }
}
