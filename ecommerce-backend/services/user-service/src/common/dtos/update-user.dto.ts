import z from 'zod';

const updateUserSchema = z.object({
  firstName: z.string().min(1, 'Name is required').optional(),
  lastName: z.string().min(1, 'Name is required').optional(),
  email: z.email('Invalid email address').optional(),
  phoneNumber: z
    .string()
    .min(10, 'Phone Number must be at least 10 digits')
    .optional(),
  password: z
    .string()
    .min(6, 'Password must be at least 6 characters')
    .optional(),
});
export type updateUserDto = z.infer<typeof updateUserSchema>;
