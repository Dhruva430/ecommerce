import {
  ConflictException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { type LoginDto } from '@ecommerce-backend/shared-dtos';
import { type SignupDto } from '@ecommerce-backend/shared-dtos';
import * as jwt from 'jsonwebtoken';
import { Role, SellerStatus } from '@prisma/client';

import { TokenUtil } from '../common/utils/token';
import { prisma } from '../common/config/prisma';
@Injectable()
export class AppService {
  constructor(private tokenUtil: TokenUtil) {}
  private otps = new Map<string, string>();

  async signup(dto: SignupDto) {
    const existingUser = await prisma.user.findFirst({
      where: { email: dto.email },
    });
    if (existingUser) {
      throw new ConflictException('Email already in use');
    }

    const user = await prisma.user.create({
      data: {
        email: dto.email,
        role: dto.role,
        credentials: {
          create: {
            password: dto.password,
          },
        },
      },
    });
    switch (dto.role) {
      case Role.SELLER:
        await prisma.seller.create({
          data: {
            userId: user.id,
            status: SellerStatus.PENDING,
          },
        });
        break;
      case Role.ADMIN:
        await prisma.admin.create({
          data: {
            userId: user.id,
          },
        });
        break;

      case Role.BUYER:
        await prisma.buyer.create({
          data: {
            userId: user.id,
          },
        });
        break;
      default:
        console.log('No role matched');
        break;
    }

    const token = this.tokenUtil.generateTokenPair({
      id: user.id,
      role: user.role,
    });
    return { message: 'Signup successful', token: token };
  }

  async login(dto: LoginDto) {
    const user = await prisma.user.findFirst({
      where: { email: dto.email, role: dto.role },
      include: { credentials: true },
    });
    if (!user || user.credentials!.password !== dto.password) {
      throw new NotFoundException('User does not exist');
    }
    const token = this.tokenUtil.generateTokenPair({
      id: user.id,
      role: user.role,
    });
    return { message: 'Login successful', token: token };
  }

  async me(headers: Record<string, string>) {
    const userId = headers['x-user-id'];
    if (!userId) {
      throw new NotFoundException('UserId not found');
    }
    const role = headers['x-user-role'] as Role;
    const user = prisma.user.findUnique({
      where: { id: userId, role: role },
      select: {
        id: true,
        role: true,
      },
    });
    return user;
  }
  async logout(headers: Record<string, string>) {
    const userId = headers['x-user-id'];
    if (!userId) {
      throw new NotFoundException('UserId not found');
    }
    await prisma.refreshToken.deleteMany({
      where: { userId: userId },
    });

    return { message: 'Logged out successfully' };
  }

  // requestOtp(dto: RequestOtpDto) {
  //   const otp = Math.floor(100000 + Math.random() * 900000).toString();
  //   this.otps.set(dto.email, otp);

  //   return { message: 'OTP generated', otp };
  // }

  verifyOtp(email: string, otp: string) {
    const stored = this.otps.get(email);
    if (!stored || stored !== otp) {
      throw new NotFoundException('User does not exist');
    }

    this.otps.delete(email);

    const token = jwt.sign(
      { id: 'user123', email: email, role: 'USER' },
      process.env.JWT_SECRET!,
      { expiresIn: '7d' }
    );

    return { accessToken: token };
  }
  async refresh(headers: Record<string, string>, res: Response) {
    const rawRefreshToken = headers['x-refresh-token'];
    if (!rawRefreshToken) {
      throw new NotFoundException('Refresh token not found');
    }
    const tokens = await this.tokenUtil.refresh(rawRefreshToken);
    res.headers.set('x-access-token', tokens.accessToken);
    res.headers.set('x-refresh-token', await tokens.refreshToken);
    return { message: 'Tokens refreshed' };
  }
  ping() {
    return { message: 'Auth service is alive' };
  }
}
