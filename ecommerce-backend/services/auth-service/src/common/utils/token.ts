import { prisma } from '@ecommerce-backend/shared';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';
export class TokenUtil {
  constructor(private jwt: JwtService) {}

  generateAccessToken(payload: { id: string; email: string; role: string }) {
    return this.jwt.sign(
      {
        sub: payload.id,
        email: payload.email,
        role: payload.role,
      },
      { expiresIn: '15m' }
    );
  }

  generateRefreshToken(payload: { id: string }) {
    return this.jwt.sign(
      {
        sub: payload.id,
      },
      { expiresIn: '7d' }
    );
  }

  async storeRefreshToken(userId: string, refreshToken: string) {
    const hashed = await this.hashed(refreshToken);
    await prisma.refreshToken.upsert({
      where: { userId: userId },
      update: { token: hashed },
      create: {
        userId: userId,
        token: hashed,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      },
    });
  }
  async generateTokenPair(user: { id: string; email: string; role: string }) {
    const accessToken = this.generateAccessToken(user);
    const refreshToken = this.generateRefreshToken({ id: user.id });
    this.storeRefreshToken(user.id, refreshToken);
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
