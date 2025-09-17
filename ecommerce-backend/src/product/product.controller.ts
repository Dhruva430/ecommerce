import { Controller, Get, Param } from '@nestjs/common';
import { ProductService } from './product.service';
@Controller('/product')
export class ProductController {
  constructor(private readonly productService: ProductService) {}
  @Get('/:id')
  getProduct(@Param('id') id: string): string {
    return this.productService.getProduct(id);
  }
}
