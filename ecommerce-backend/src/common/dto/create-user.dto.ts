import { Role } from '@prisma/client';
import z from 'zod';

export const CreateUserSchema = z.object({
  email: z.email(),
  password: z.string().min(6).max(100),
  role: z.enum(Role).default(Role.USER),
});

export type CreateUserDto = z.infer<typeof CreateUserSchema>;
