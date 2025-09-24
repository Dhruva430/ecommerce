import { Injectable } from '@nestjs/common';
import { prisma } from 'lib/primsa';
import { type CreateProductDto } from './dto/product.dto';

@Injectable()
export class ProductService {
  async getProduct(id: number) {
    const product = await prisma.product.findUnique({
      where: { id },
    });
    return product;
  }
  async createProduct(createProductDto: CreateProductDto) {
    await prisma.product.create({
      data: createProductDto,
    });
    return {
      message: 'Product created successfully',
      product: createProductDto,
    };
  }
  async deleteProduct(id: number) {
    await prisma.product.delete({
      where: { id },
    });
    return {
      message: 'Product deleted successfully',
    };
  }
}
