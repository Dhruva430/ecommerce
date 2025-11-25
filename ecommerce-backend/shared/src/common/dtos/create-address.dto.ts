import z from 'zod';

export const CreateAddressSchema = z.object({
  pincode: z.string().min(5).max(10),
  city: z.string().min(2).max(50),
  state: z.string().min(2).max(50),
  landmark: z.string().min(2).max(100).optional(),
  area: z.string().min(2).max(100),
  houseNo: z.string().min(1).max(20),
  phoneNumber: z.string().min(10).max(15),
  name: z.string().min(2).max(50),
});
export type CreateAddressDto = z.infer<typeof CreateAddressSchema>;
