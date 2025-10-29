import { Category } from '@prisma/client';
import z from 'zod';

const VariantImageSchema = z.object({
  imageUrl: z.url(),
  position: z.number().optional(),
});

const VariantAttributeSchema = z.object({
  name: z.string(),
  value: z.string(),
});

const ProductVariantSchema = z.object({
  title: z.string(),
  description: z.string().optional(),
  price: z.number().min(0),
  stock: z.number().int().min(0),
  size: z.string().optional(),
  variantImages: z.array(VariantImageSchema).optional(),
  variantAttributes: z.array(VariantAttributeSchema).optional(),
});

export const CreateProductSchema = z.object({
  title: z.string().min(2),
  description: z.string().optional(),
  price: z.number().min(0),
  discounted: z.number().optional(),
  imageUrl: z.url(),
  category: z.enum(Category),
  variants: z.array(ProductVariantSchema).optional(),
});
export type CreateProductDto = z.infer<typeof CreateProductSchema>;
