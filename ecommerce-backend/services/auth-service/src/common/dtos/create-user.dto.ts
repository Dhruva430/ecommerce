import { Role } from '@prisma/client';
import z from 'zod';

export const createUserSchema = z.object({
  name: z.string().min(2),
  email: z.email(),
  role: z.enum(Role).default(Role.BUYER),
  password: z.string().min(6),
});

export type CreateUserDto = z.infer<typeof createUserSchema>;
