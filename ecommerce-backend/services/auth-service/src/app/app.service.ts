import { Injectable, UnauthorizedException } from '@nestjs/common';

import { LoginSchema, LoginDto } from '../common/dtos/login';
import { SignupSchema, SignupDto } from '../common/dtos/signup';
import { RequestOtpSchema, RequestOtpDto } from '../common/dtos/request-otp';
import * as jwt from 'jsonwebtoken';
import { Role } from '@prisma/client';

@Injectable()
export class AppService {
  private otps = new Map<string, string>();

  signup(dto: SignupDto) {
    // replace with DB check later
    return { message: 'User registered' };
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

    return { message: 'OTP generated', otp }; // remove otp in production
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
