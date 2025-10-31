import { Category } from '@prisma/client';
import { prisma } from 'libs/primsa';
import { CreateProductDto } from 'src/common/dto/create-product.dto';

export default class SellerService {
  getProducts(sellerId: string) {
    const products = prisma.product.findMany({
      where: { sellerId: sellerId },
      include: {
        productVariant: {
          include: { variantImages: true, variantAttribute: true },
        },
      },
    });
    return products;
  }

  async createProduct(userId: string, dto: CreateProductDto) {
    const seller = await prisma.seller.findUnique({
      where: {
        userId,
      },
    });
    if (!seller) {
      throw new Error('Seller not found');
    }
    const product = await prisma.product.create({
      data: {
        title: dto.title,
        description: dto.description || '',
        price: dto.price,
        imageUrl: dto.imageUrl,
        category: Category.ELECTRONICS,
        discounted: dto.discounted || 0,
        sellerId: seller.id,
        productVariant: {
          create: dto.variants?.map((variant) => ({
            title: variant.title,
            description: variant.description || '',
            price: variant.price,
            stock: variant.stock,
            size: variant.size || '',

            variantAttribute: variant.variantAttributes
              ? {
                  create: variant.variantAttributes.map((attr) => ({
                    name: attr.name,
                    value: attr.value,
                  })),
                }
              : undefined,
            variantImages: variant.variantImages
              ? {
                  create: variant.variantImages.map((img) => ({
                    imageUrl: img.imageUrl,
                    position: img.position,
                  })),
                }
              : undefined,
          })),
        },
      },
      include: {
        productVariant: {
          include: { variantImages: true, variantAttribute: true },
        },
      },
    });
    return product;
  }

  async deleteProduct(productId: string) {
    const deletedProduct = await prisma.product.deleteMany({
      where: { id: productId },
    });
    return deletedProduct;
  }

  async getSellerById(sellerId: string) {
    return prisma.seller.findUnique({
      where: { id: sellerId },
    });
  }

  async updateProduct(sellerId: string, productId: string, dto: any) {
    const updatedProduct = await prisma.product.updateMany({
      where: { id: productId, sellerId: sellerId },
      data: {
        title: dto.title,
        description: dto.description,
        price: dto.price,
        imageUrl: dto.imageUrl,
        category: dto.category,
        discounted: dto.discounted,
      },
    });
    return updatedProduct;
  }
  async getOrders(sellerId: string) {
    const orders = await prisma.order.findMany({
      where: { sellerId: sellerId },
    });
    return orders;
  }
}
