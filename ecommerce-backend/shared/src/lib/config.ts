import z from 'zod';

const configSchema = z.object({
  DATABASE_URL: z.url(),
  JWT_SECRET: z.string().min(10),
});
export const config: z.infer<typeof configSchema> = configSchema.parse(
  process.env
);
