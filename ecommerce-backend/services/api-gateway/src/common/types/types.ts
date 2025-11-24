import { Role } from '@prisma/client';

export type JwtPayload = {
  id: string;
  role: Role;
  ip?: string;
  device?: string;
};
