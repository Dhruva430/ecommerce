import { Controller, Get, Query } from '@nestjs/common';
import { CategoryService } from './category.service';

@Controller('/categories')
export class CategoryController {
  constructor(private readonly categoryService: CategoryService) {}

  @Get()
  findAll(): string {
    return this.categoryService.findAll();
  }
  @Get()
  findById(@Query('name') name: string): string {
    return this.categoryService.findById(name);
  }
}
