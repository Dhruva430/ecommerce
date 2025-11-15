import {
  Injectable,
  NestMiddleware,
  UnauthorizedException,
} from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';

@Injectable()
export class AuthMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const token = req.headers['authorization']?.split(' ')[1];
    if (!token) throw new UnauthorizedException('Missing token');

    try {
      const payload = jwt.verify(token, process.env.JWT_SECRET!);
      (req as any).user = payload;
      return next();
    } catch {
      (req as any).jwtError = true;
      return next();
    }
  }
}
