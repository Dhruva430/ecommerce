import { Injectable } from '@nestjs/common';
import { createProxyMiddleware } from 'http-proxy-middleware';
import { Application } from 'express';

@Injectable()
export class ProxyService {
  register(app: Application) {
    app.use(
      '/auth',
      createProxyMiddleware({
        target: process.env.AUTH_SERVICE_URL,
        changeOrigin: true,
        pathRewrite: { '^/auth': '' },
        on: {
          proxyReq: (proxyReq, req, res) => {
            // 🔥 Forward Authorization header (required for your AuthGuard)
            const authHeader = req.headers['authorization'];
            if (authHeader) {
              proxyReq.setHeader('authorization', authHeader);
            }

            // forward other custom user headers if available
            const user = (req as any).user;
            if (user) {
              proxyReq.setHeader('x-user-id', user.id);
              proxyReq.setHeader('x-user-role', user.role);
            }
          },
        },
      })
    );
  }
}
