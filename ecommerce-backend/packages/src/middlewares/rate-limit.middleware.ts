import rateLimit from 'express-rate-limit';
import { RequestHandler } from 'express';

export const RateLimitMiddleware: RequestHandler = rateLimit({
  windowMs: Number(process.env.RATE_LIMIT_WINDOW_MS || 60_000),
  max: Number(process.env.RATE_LIMIT_MAX || 100),
  standardHeaders: true,
  legacyHeaders: false,
  handler: (req, res) => {
    res.status(429).json({ success: false, message: 'Too many requests' });
  },
}) as unknown as RequestHandler;
