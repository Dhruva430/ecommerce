import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  register(userData: any) {
    // Registration logic here
    return { message: 'User registered successfully', userData };
  }
}
