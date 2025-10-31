import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  ParseIntPipe,
  Post,
  Query,
  UsePipes,
} from '@nestjs/common';
import { ProductService } from './product.service';

@Controller('/api/products')
export class ProductController {
  constructor(private readonly productService: ProductService) {}
  @Get()
  getAllProducts(
    @Query('limit', ParseIntPipe) limit: number,
    @Query('offset', ParseIntPipe) offset: number,
  ) {
    return this.productService.getProducts(limit, offset);
  }
  @Get('/restore')
  restoreProducts() {
    return this.productService.restoreProducts();
  }

  @Get('/:productId')
  getProductById(@Param('productId') productId: string) {
    return this.productService.getProductById(productId);
  }
  @Get('/category/:category')
  getProductsByCategory(@Param('category') category: string) {
    return this.productService.getProductByCategory(category);
  }
}
