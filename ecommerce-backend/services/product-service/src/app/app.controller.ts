import {
  Body,
  Controller,
  Get,
  ParseIntPipe,
  Post,
  Query,
} from '@nestjs/common';
import { AppService } from './app.service';

@Controller('/products')
export class AppController {
  constructor(private readonly productService: AppService) {}

  @Get()
  getAllProducts(
    @Query('limit', ParseIntPipe) limit: number,
    @Query('cursor') cursor?: string
  ) {
    return this.productService.getAllProducts(limit, cursor);
  }
}
