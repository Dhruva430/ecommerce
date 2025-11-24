import { z } from 'zod';

export const SignupSchema = z.object({
  email: z.email(),
  role: z.enum(['BUYER', 'SELLER', 'ADMIN']),
  password: z.string().min(6),
});

export type SignupDto = z.infer<typeof SignupSchema>;
