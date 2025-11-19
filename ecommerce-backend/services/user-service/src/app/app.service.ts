import { cleanDto, prisma } from '@ecommerce-backend/shared';
import { Injectable, NotFoundException } from '@nestjs/common';
import { type createAddressDto } from '../common/dtos/create-address.dto';
import { type updateAddressDto } from '../common/dtos/update-address.dto';
import { type updateUserDto } from '../common/dtos/update-user.dto';
@Injectable()
export class AppService {
  // ------------------- User Management -------------------
  async getUser(userId: string) {
    const user = await prisma.user.findUnique({
      where: { id: userId },
      include: {
        order: true,
      },
    });
    return user;
  }
  async deleteUser(req: Request) {
    const userId = req.headers.get('x-user-id') as string;
    await prisma.user.delete({
      where: { id: userId },
    });
    return { message: 'User deleted successfully' };
  }
  async updateUser(req: Request, dto: updateUserDto) {
    const userId = req.headers.get('x-user-id') as string;
    const data = cleanDto(dto);

    await prisma.user.update({
      where: { id: userId },
      data: { ...data },
    });
    return { message: 'User updated successfully' };
  }
  // ------------------ Address Management ------------------
  async getUserAddress(req: Request) {
    const userId = req.headers.get('x-user-id') as string;
    const addresses = await prisma.address.findMany({
      where: { userId: userId },
    });
    return addresses;
  }
  async addUserAddress(req: Request, dto: createAddressDto) {
    const userId = req.headers.get('x-user-id') as string;
    if (!userId) {
      throw new NotFoundException('User ID not found in request headers');
    }
    const address = await prisma.address.create({
      data: { ...dto, userId: userId },
    });
    return address;
  }
  async deleteUserAddress(req: Request, addressId: string) {
    const userId = req.headers.get('x-user-id') as string;
    await prisma.address.deleteMany({
      where: { userId: userId, id: addressId },
    });
    return { message: 'User address deleted successfully' };
  }
  async updateUserAddress(
    addressId: string,
    dto: updateAddressDto,
    req: Request
  ) {
    const userId = req.headers.get('x-user-id') as string;
    const data = cleanDto(dto);
    await prisma.address.update({
      where: { id: addressId, userId: userId },
      data: { ...data },
    });
    return { message: 'User address updated successfully' };
  }

  // ------------------- User Activity -------------------
  async getUserOrders(req: Request) {
    const userId = req.headers.get('x-user-id') as string;
    const orders = await prisma.order.findMany({
      where: { userId: userId },
      include: {
        orderProduct: true,
      },
    });
    return orders;
  }
  async getUserReviews(req: Request) {
    const userId = req.headers.get('x-user-id') as string;
    const reviews = await prisma.reviews.findMany({
      where: { userId: userId },
    });
    return reviews;
  }
}
