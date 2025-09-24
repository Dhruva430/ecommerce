import { Body, Controller, Delete, Get, Param, Post } from '@nestjs/common';
import { ProductService } from './product.service';
import { type CreateProductDto } from './dto/product.dto';
@Controller('/product')
export class ProductController {
  constructor(private readonly productService: ProductService) {}
  @Get('/:id')
  getProduct(@Param('id') id: number) {
    return this.productService.getProduct(id);
  }
  @Post('/create')
  async createProduct(@Body() createProductDto: CreateProductDto) {
    return await this.productService.createProduct(createProductDto);
  }
  @Delete('/:id')
  deleteProduct(@Param('id') id: number) {
    return this.productService.deleteProduct(id);
  }
}
