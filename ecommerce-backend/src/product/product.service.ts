import { Injectable } from '@nestjs/common';
import productData from '../data/products.json';
@Injectable()
export class ProductService {
  private products = productData;
  getAllProducts() {
    return this.products;
  }
  getProductById(id: string) {
    return this.products.find((product) => product.id === id);
  }
  getProductByCategory(category: string) {
    const categories = this.products.filter(
      (product) =>
        product.category &&
        product.category.toLowerCase() === category.toLowerCase(),
    );

    return categories;
  }
}
