import { packages } from '@ecommerce-backend/shared';
import { Injectable } from '@nestjs/common';
@Injectable()
export class AppService {
  getData(): { message: string } {
    const secret = packages();
    console.log('Shared secret:', secret);
    return { message: 'Hello API' };
  }
}
