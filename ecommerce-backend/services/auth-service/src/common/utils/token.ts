import { Role } from '@prisma/client';
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { prisma } from '../config/prisma';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcryptjs';
import { JwtPayload } from '../../../../../shared/src/common/types/types';
import { config } from 'shared/src/lib/config';

@Injectable()
export class TokenUtil {
  constructor(private readonly jwtService: JwtService) {}

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
    const refreshToken = this.generateRefreshToken({
      id: userId,
      role,
      device: 'unknown',
      ip: 'unknown',
    });

    const hashed = await this.hashed(refreshToken);

    await prisma.refreshToken.upsert({
      where: { userId },
      update: {
        token: hashed,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      },
      create: {
        userId,
        token: hashed,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      },
    });

    return refreshToken;
  }

  async refresh(rawRefreshToken: string) {
    try {
      const payload: JwtPayload = this.jwtService.verify(rawRefreshToken, {
        secret: config.JWT_SECRET,
      });
      const storedTokens = await prisma.refreshToken.findMany({
        where: { userId: payload.id },
      });

      if (!storedTokens.length) {
        throw new UnauthorizedException('No refresh token found');
      }

      const matched = storedTokens.find((t) =>
        bcrypt.compareSync(rawRefreshToken, t.token)
      );

      if (!matched) throw new UnauthorizedException('Invalid refresh token');

      if (new Date() > matched.expiresAt) {
        throw new UnauthorizedException('Refresh token expired');
      }

      await prisma.refreshToken.update({
        where: { id: matched.id },
        data: { revoked: true },
      });

      return this.generateTokenPair(payload);
    } catch (e) {
      throw new UnauthorizedException('Failed to refresh token');
    }
  }

  async generateTokenPair(user: JwtPayload) {
    return {
      accessToken: this.generateAccessToken(user),
      refreshToken: await this.createAndStoreRefreshToken(
        user.id,
        user.role as Role
      ),
    };
  }

  hashed(data: string) {
    return bcrypt.hash(data, 10);
  }
}
