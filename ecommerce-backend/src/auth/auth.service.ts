import { Injectable, UnauthorizedException } from '@nestjs/common';

import { JwtService } from '@nestjs/jwt';
import { prisma } from 'libs/primsa';
import { type CreateUserDto } from 'src/common/dto/create-user.dto';
import hashPassword, { comparePassword } from 'src/common/utils';

@Injectable()
export default class AuthService {
  constructor(private jwtService: JwtService) {}

  async register(dto: CreateUserDto) {
    const existingUser = await prisma.user.findFirst({
      where: { email: dto.email, role: dto.role },
    });
    if (existingUser) {
      throw new UnauthorizedException('User already exists');
    }
    const hashedPassword = hashPassword(dto.password);
    const user = await prisma.user.create({
      data: {
        email: dto.email,
        role: dto.role,
        credentials: {
          create: {
            password: hashedPassword,
          },
        },
      },
    });

    const token = await this.jwtService.signAsync({
      id: user.id,
      email: user.email,
      role: user.role,
    });
    return { user, token };
  }
  async login(dto: CreateUserDto) {
    const user = await prisma.user.findFirst({
      where: { email: dto.email, role: dto.role },
      include: { credentials: true },
    });
    if (!user) {
      throw new UnauthorizedException('Email is not registered');
    }
    if (comparePassword(dto.password, user.credentials!.password) === false) {
      throw new UnauthorizedException('Invalid password');
    }
    const token = await this.jwtService.signAsync({
      id: user.id,
      email: user.email,
      role: user.role,
    });
    return { user, token };
  }
}
