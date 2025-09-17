import { Injectable } from '@nestjs/common';

@Injectable()
export class ProductService {
  getProduct(id: string): string {
    return `This action returns product with id: ${id}`;
  }
}
