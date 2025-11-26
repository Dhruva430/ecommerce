import { Injectable } from '@nestjs/common';
import { PrismaService } from '@ecommerce-backend/shared';
@Injectable()
export class AppService {
  constructor(private prisma: PrismaService) {}
  async getData() {
    const data = await this.prisma.user.findMany();
    return data;
  }
}
