import { z } from 'zod';

export const LoginSchema = z.object({
  email: z.email(),
  role: z.enum(['BUYER', 'SELLER', 'ADMIN']),
  password: z.string().min(6),
});

export type LoginDto = z.infer<typeof LoginSchema>;
