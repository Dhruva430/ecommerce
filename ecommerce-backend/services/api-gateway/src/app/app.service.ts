import { Injectable } from '@nestjs/common';
import { packages } from '@ecommerce-backend/shared';

@Injectable()
export class AppService {
  getData(): { message: string } {
    const data = packages();
    console.log('Shared module says:', data);
    return { message: 'Hello API' };
  }
}
