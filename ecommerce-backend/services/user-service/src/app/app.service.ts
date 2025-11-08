import { shared } from '@ecommerce-backend/shared';
import { Injectable } from '@nestjs/common';
@Injectable()
export class AppService {
  getData(): { message: string } {
    const secret = shared();
    console.log('Shared secret:', secret);
    return { message: 'Hello API' };
  }
}
