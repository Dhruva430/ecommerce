import { Injectable } from '@nestjs/common';
import { createProxyMiddleware } from 'http-proxy-middleware';
import { Application, Request, Response } from 'express';

@Injectable()
export class ProxyService {
  private forwardRequest(proxyReq: any, req: Request) {
    const authHeader = req.headers['authorization'];
    if (authHeader) {
      proxyReq.setHeader('authorization', authHeader);
    }
    if (req.user) {
      proxyReq.setHeader('x-user-id', req.user.id);
      proxyReq.setHeader('x-user-role', req.user.role);
      proxyReq.setHeader('x-user-email', req.user.email);
    }

    // 🔥 Forward JSON body for POST/PUT/PATCH
    if (req.body && Object.keys(req.body).length) {
      const bodyData = JSON.stringify(req.body);
      proxyReq.setHeader('Content-Type', 'application/json');
      proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
      proxyReq.write(bodyData);
    }
  }
  private createServiceProxy(targetUrl: string) {
    return createProxyMiddleware({
      target: targetUrl,
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
      on: {
        proxyReq: (proxyReq, req: Request, res: Response) => {
          this.forwardRequest(proxyReq, req);
        },
        error: (err, req, res) => {
          console.error('Proxy error:', err);
          res.end('Internal proxy error occurred. Please try again later.');
        },
      },
    });
  }
  register(app: Application) {
    app.use(
      '/api/auth',
      this.createServiceProxy(process.env.AUTH_SERVICE_URL!)
    );
    app.use(
      '/api/users',
      this.createServiceProxy(process.env.USER_SERVICE_URL!)
    );
  }
}
