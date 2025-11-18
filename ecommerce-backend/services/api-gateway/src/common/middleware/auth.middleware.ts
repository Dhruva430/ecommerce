import {
  Injectable,
  NestMiddleware,
  UnauthorizedException,
} from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import jwt, { JwtPayload } from 'jsonwebtoken';

@Injectable()
export class AuthMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const authHeader = req.headers['authorization'];
    if (!authHeader) {
      throw new UnauthorizedException('Missing Authorization header');
    }
    const token = authHeader.split(' ')[1];

    try {
      const payload = jwt.verify(
        token,
        process.env.JWT_SECRET!
      ) as JwtPayload & {
        id: string;
        role: string;
        email?: string;
      };

      req.user = {
        id: payload.id,
        role: payload.role,
        email: payload.email,
      };

      next();
    } catch (err: unknown) {
      throw new UnauthorizedException('Invalid token');
    }
  }
}
