import { prisma } from 'libs/primsa';

export class UserService {
  async getProfile(userId: string) {
    const user = await prisma.user.findUnique({
      where: { id: userId },
    });
    return user;
  }
  async getOrders(userId: string) {
    const orders = await prisma.order.findMany({
      where: { userId },
    });
    return orders;
  }
  async updateProfile(userId: string, dto: any) {
    const updatedUser = await prisma.user.update({
      where: { id: userId },
      data: dto,
    });
    return updatedUser;
  }
}
