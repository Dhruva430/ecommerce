import { Injectable, UnauthorizedException } from '@nestjs/common';

import { type LoginDto } from '../common/dtos/login';
import { type SignupDto } from '../common/dtos/signup';
import { type RequestOtpDto } from '../common/dtos/request-otp';
import * as jwt from 'jsonwebtoken';
import { Role, SellerStatus } from '@prisma/client';
import { prisma } from '@ecommerce-backend/shared';

@Injectable()
export class AppService {
  private otps = new Map<string, string>();

  async signup(dto: SignupDto) {
    const existingUser = await prisma.user.findFirst({
      where: { email: dto.email },
    });
    if (existingUser) {
      throw new UnauthorizedException('Email already in use');
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
    const token = jwt.sign(
      { id: user.id, email: user.email, role: user.role },
      process.env.JWT_SECRET!,
      { expiresIn: '7d' }
    );
    return { message: 'Signup successful', token: token };
  }

  login(dto: LoginDto) {
    // Replace with DB auth logic
    if (dto.email !== 'test@test.com' || dto.password !== '123456')
      throw new UnauthorizedException('Invalid credentials');

    const token = jwt.sign(
      { id: 'user123', email: dto.email, role: Role.BUYER },
      process.env.JWT_SECRET!,
      { expiresIn: '7d' }
    );

    return { accessToken: token };
  }

  requestOtp(dto: RequestOtpDto) {
    const otp = Math.floor(100000 + Math.random() * 900000).toString();
    this.otps.set(dto.email, otp);

    return { message: 'OTP generated', otp };
  }

  verifyOtp(email: string, otp: string) {
    const stored = this.otps.get(email);
    if (!stored || stored !== otp) {
      throw new UnauthorizedException('Invalid OTP');
    }

    this.otps.delete(email);

    const token = jwt.sign(
      { id: 'user123', email: email, role: 'USER' },
      process.env.JWT_SECRET!,
      { expiresIn: '7d' }
    );

    return { accessToken: token };
  }
}
