import z from 'zod';

export const updateUserSchema = z.object({
  firstName: z.string().min(2).max(50).optional(),
  lastName: z.string().min(2).max(50).optional(),
  email: z.email().optional(),
  phoneNumber: z.string().min(10).max(15).optional(),
  password: z.string().min(6).max(100).optional(),
});

export type UpdateUserDto = z.infer<typeof updateUserSchema>;
