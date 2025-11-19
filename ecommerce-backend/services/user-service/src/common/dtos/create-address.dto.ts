import { z } from 'zod';
const createAddressSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  city: z.string().min(1, 'City is required'),
  state: z.string().min(1, 'State is required'),
  pincode: z.number().min(1, 'Zip Code is required'),
  area: z.string().min(1, 'Area is required'),
  landmark: z.string().optional(),
  phoneNumber: z.number().min(1, 'Phone Number is required'),
  houseNo: z.string().min(1, 'House Number is required'),
});
export type createAddressDto = z.infer<typeof createAddressSchema>;
