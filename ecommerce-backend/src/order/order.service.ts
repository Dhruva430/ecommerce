import { BadRequestException, Injectable } from '@nestjs/common';
import { prisma } from 'libs/primsa';
import { RedlockService } from 'libs/redlock';

@Injectable()
export class OrderService {
  constructor(private readonly redlock: RedlockService) {}
  async buyOrder(
    userId: string,
    variantId: string,
    quantity: number,
    address: string,
  ) {
    const lockKey = `locks:product_variant:${variantId}`;
    const lock = await this.redlock.acquireLock(lockKey, 5000);

    try {
      const variant = await prisma.productVariant.findUnique({
        where: { id: variantId },
      });

      if (!variant) {
        throw new BadRequestException('Product variant not found');
      }

      if (variant.stock < quantity) {
        throw new BadRequestException('Insufficient stock');
      }

      // Run DB ops in a transaction for safety
      const result = await prisma.$transaction(async (tx) => {
        // Decrement stock
        await tx.productVariant.update({
          where: { id: variantId },
          data: { stock: { decrement: quantity } },
        });

        // Fetch seller ID for order
        const product = await tx.product.findUnique({
          where: { id: variant.productId },
          select: { sellerId: true },
        });

        // Create the order
        const order = await tx.order.create({
          data: {
            userId,
            Address: 'User provided address',
            totalAmount: variant.price * quantity,
            status: 'PENDING',
            paymentStatus: 'PENDING',
            sellerId: product?.sellerId,
            orderProduct: {
              create: {
                productId: variant.productId,
                sellerId: product?.sellerId!,
                variantId,
                amount: quantity,
              },
            },
          },
          include: { orderProduct: true },
        });

        return order;
      });

      return { message: 'Order placed successfully', order: result };
    } catch (err) {
      throw new BadRequestException(err.message || 'Order failed');
    } finally {
      await this.redlock.releaseLock(lock);
    }
  }
}
