import { Injectable, UnauthorizedException } from '@nestjs/common';

import { JwtService } from '@nestjs/jwt';
import { prisma } from 'libs/primsa';

@Injectable()
export class AuthService {}
