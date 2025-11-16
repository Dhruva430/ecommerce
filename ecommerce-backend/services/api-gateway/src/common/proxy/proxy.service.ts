import { Injectable } from '@nestjs/common';
import { createProxyMiddleware } from 'http-proxy-middleware';
import { Application, Request, Response } from 'express';

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
          proxyReq: (proxyReq, req: Request, res: Response) => {
            const authHeader = req.headers['authorization'];
            if (authHeader) {
              proxyReq.setHeader('authorization', authHeader);
            }
            if (req.user) {
              proxyReq.setHeader('x-user-id', req.user.id);
              proxyReq.setHeader('x-user-role', req.user.role);
            }

            // 🔥 Forward JSON body for POST/PUT/PATCH
            if (req.body && Object.keys(req.body).length) {
              const bodyData = JSON.stringify(req.body);
              proxyReq.setHeader('Content-Type', 'application/json');
              proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
              proxyReq.write(bodyData);
            }
          },
        },
      })
    );
    app.use(
      '/user',
      createProxyMiddleware({
        target: process.env.USER_SERVICE_URL,
        changeOrigin: true,
        pathRewrite: { '^/user': '' },
        on: {
          proxyReq: (proxyReq, req: Request, res: Response) => {
            // 🔥 Forward Authorization header (required for your AuthGuard)
            const authHeader = req.headers['authorization'];
            if (authHeader) {
              proxyReq.setHeader('authorization', authHeader);
            }

            if (req.user) {
              proxyReq.setHeader('x-user-id', req.user.id);
              proxyReq.setHeader('x-user-role', req.user.role);
            }

            if (req.body && Object.keys(req.body).length) {
              const bodyData = JSON.stringify(req.body);
              proxyReq.setHeader('Content-Type', 'application/json');
              proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
              proxyReq.write(bodyData);
            }
          },
          error: (err, req, res) => {
            console.error('Proxy error:', err);
            res.end(
              'Something went wrong. And we are reporting a custom error message.'
            );
          },
        },
      })
    );
  }
}
