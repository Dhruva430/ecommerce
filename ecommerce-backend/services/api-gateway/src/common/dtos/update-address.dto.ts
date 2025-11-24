import z from 'zod';

export const UpdateAddressSchema = z.object({
  pincode: z.string().min(5).max(10).optional(),
  city: z.string().min(2).max(50).optional(),
  state: z.string().min(2).max(50).optional(),
  landmark: z.string().min(2).max(100).optional(),
  area: z.string().min(2).max(100).optional(),
  houseNo: z.string().min(1).max(20).optional(),
  phoneNumber: z.string().min(10).max(15).optional(),
  name: z.string().min(2).max(50).optional(),
});

export type UpdateAddressDto = z.infer<typeof UpdateAddressSchema>;
