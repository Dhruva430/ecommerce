import z from 'zod';
export const updateAddressSchema = z.object({
  name: z.string().min(1, 'Name is required').optional(),
  city: z.string().min(1, 'City is required').optional(),
  state: z.string().min(1, 'State is required').optional(),
  pincode: z.number().min(1, 'Zip Code is required').optional(),
  area: z.string().min(1, 'Area is required').optional(),
  landmark: z.string().optional(),
  phoneNumber: z.number().min(1, 'Phone Number is required').optional(),
  houseNo: z.string().min(1, 'House Number is required').optional(),
});
export type updateAddressDto = z.infer<typeof updateAddressSchema>;
