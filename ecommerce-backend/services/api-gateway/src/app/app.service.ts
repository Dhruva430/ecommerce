import { shared } from '@ecommerce-backend/shared';
import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getData(): { message: string } {
    const data = shared();
    console.log('Shared module says:', data);
    return { message: 'Hello API' };
  }
}
