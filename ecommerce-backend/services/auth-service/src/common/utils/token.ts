import { prisma } from '@ecommerce-backend/shared';
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { Role } from '@prisma/client';
import bcrypt from 'bcryptjs';

import { JwtPayload } from '../types/types';

@Injectable()
export class TokenUtil {
  constructor(private jwtService: JwtService) {}

  generateAccessToken(payload: JwtPayload) {
    return this.jwtService.sign(
      {
        sub: payload.id,
        role: payload.role,
        device: payload.device,
        ip: payload.ip,
      },
      { expiresIn: '15m' }
    );
  }

  generateRefreshToken(payload: JwtPayload) {
    return this.jwtService.sign(
      {
        sub: payload.id,
        role: payload.role,
        device: payload.device,
        ip: payload.ip,
      },
      { expiresIn: '7d' }
    );
  }

  async createAndStoreRefreshToken(userId: string, role: Role) {
    const device = 'unknown'; // TODO: get device info from request
    const ip = 'unknown'; // TODO: get ip address from request
    const refreshToken = this.generateRefreshToken({
      id: userId,
      role,
      device,
      ip,
    });
    const hashed = await this.hashed(refreshToken);
    await prisma.refreshToken.upsert({
      where: { userId: userId },
      update: {
        token: hashed,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      }, // TODO: add device and ip address after ngrok testing
      create: {
        userId: userId,
        token: hashed,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      },
    });
    return refreshToken;
  }

  async refresh(rawRefreshToken: string) {
    try {
      const payload: JwtPayload = this.jwtService.verify(rawRefreshToken, {
        secret: process.env.JWT_SECRET,
      });
      const userId = payload.id;
      const token = await prisma.refreshToken.findMany({
        where: { userId: userId },
      });
      if (token.length === 0 && !token) {
        throw new UnauthorizedException('No refresh token found');
      }
      let matched = null;
      for (const t of token) {
        const isMatch = this.compareHash(rawRefreshToken, t.token);
        if (isMatch) {
          matched = t;
          break;
        }
      }
      if (!matched) {
        throw new UnauthorizedException('Invalid refresh token');
      }
      if (new Date() > matched.expiresAt) {
        await prisma.refreshToken.update({
          where: { id: matched.id },
          data: { revoked: true },
        });
        throw new UnauthorizedException('Refresh token expired');
      }
      await prisma.refreshToken.update({
        where: { id: matched.id },
        data: { revoked: true },
      });
      const newRefreshToken = this.createAndStoreRefreshToken(
        payload.id,
        payload.role
      );
      const newAccessToken = this.generateAccessToken({
        id: payload.id,
        role: payload.role,
        device: payload.device,
        ip: payload.ip,
      });

      return { accessToken: newAccessToken, refreshToken: newRefreshToken };
    } catch (e) {
      throw new Error('Error while refreshing token: ');
    }
  }
  async generateTokenPair(user: JwtPayload) {
    const accessToken = this.generateAccessToken(user);

    const refreshToken = await this.createAndStoreRefreshToken(
      user.id,
      user.role as Role
    );

    return { accessToken, refreshToken };
  }
  hashed(data: string) {
    const hashed = bcrypt.hash(data, 10);
    return hashed;
  }
  compareHash(data: string, hashedData: string): boolean {
    const isMatch = bcrypt.compareSync(data, hashedData);
    return isMatch;
  }
}
