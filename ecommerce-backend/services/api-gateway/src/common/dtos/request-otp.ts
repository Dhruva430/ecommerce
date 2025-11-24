import { z } from 'zod';

export const RequestOtpSchema = z.object({
  email: z.email(),
});

export type RequestOtpDto = z.infer<typeof RequestOtpSchema>;
