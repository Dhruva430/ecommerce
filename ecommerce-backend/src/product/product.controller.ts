import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  UsePipes,
} from '@nestjs/common';
import { ProductService } from './product.service';
import { createProductDto, type CreateProductDto } from './dto/product.dto';
import { ZodValidationPipe } from 'src/pipes/zod-validation-pipe';
@Controller('/product')
export class ProductController {
  constructor(private readonly productService: ProductService) {}
  @Get('/:id')
  getProduct(@Param('id') id: string) {
    return this.productService.getProduct(id);
  }
  @Post('/create-product')
  @UsePipes(new ZodValidationPipe(createProductDto))
  async createProduct(@Body() createProductDto: CreateProductDto) {
    return await this.productService.createProduct(createProductDto);
  }
  @Delete('/delete-product/:id')
  deleteProduct(@Param('id') id: string) {
    return this.productService.deleteProduct(id);
  }
}
