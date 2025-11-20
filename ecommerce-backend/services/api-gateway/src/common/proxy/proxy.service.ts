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
    }

    // 🔥 Forward JSON body for POST/PUT/PATCH
    if (
      req.body &&
      Object.keys(req.body).length &&
      ['POST', 'PUT', 'PATCH'].includes(req.method)
    ) {
      const bodyData = JSON.stringify(req.body);
      proxyReq.setHeader('Content-Type', 'application/json');
      proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
      proxyReq.write(bodyData);
    }
  }
  private createServiceProxy(
    targetUrl: string,
    mountPath: string,
    servicePath = ''
  ) {
    const rewriteRule: Record<string, string> = {};
    rewriteRule[`^${mountPath}`] = servicePath || '';

    return createProxyMiddleware({
      target: targetUrl,
      changeOrigin: true,
      // logger: console,
      pathRewrite: rewriteRule,
      on: {
        proxyReq: (proxyReq, req: Request, res: Response) => {
          console.log(
            `[PROXY] → ${req.method} ${req.originalUrl} -> ${proxyReq.getHeader(
              'host'
            )}`
          );
          this.forwardRequest(proxyReq, req);
        },
        proxyRes: (proxyRes, req: Request, res: Response) => {
          console.log(
            `[PROXY] ← ${req.method} ${req.originalUrl} (${proxyRes.statusCode})`
          );
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
      this.createServiceProxy(process.env.AUTH_SERVICE_URL!, '/', '/auth/')
    );
    app.use(
      '/api/user',
      this.createServiceProxy(process.env.USER_SERVICE_URL!, '/', '/user/')
    );
    app.use(
      '/api/products',
      this.createServiceProxy(
        process.env.PRODUCT_SERVICE_URL!,
        '/',
        '/products/'
      )
    );
  }
}
