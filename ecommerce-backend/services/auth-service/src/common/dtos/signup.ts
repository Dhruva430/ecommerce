import { z } from 'zod';

export const SignupSchema = z.object({
  email: z.email(),
  password: z.string().min(6),
  name: z.string().min(2),
});

export type SignupDto = z.infer<typeof SignupSchema>;
