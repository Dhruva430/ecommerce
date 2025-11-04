import { Injectable } from '@nestjs/common';

@Injectable()
export class GatewayService {
  getStatus() {
    return 'Gateway operational';
  }
}
