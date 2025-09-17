import { Injectable } from '@nestjs/common';

@Injectable()
export class CategoryService {
  findAll(): string {
    return 'This action returns all categories';
  }
  findById(id: string): string {
    return `This action returns category with id: ${id}`;
  }
}
