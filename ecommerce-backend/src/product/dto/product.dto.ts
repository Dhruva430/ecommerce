import z from 'zod';
export const createProductDto = z.object({
  name: z.string().min(1),
  description: z.string().min(1),
  price: z.number().min(0),
  imageUrl: z.url(),
});
export type CreateProductDto = z.infer<typeof createProductDto>;
