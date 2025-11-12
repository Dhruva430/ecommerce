import { Body, Controller, Post, UsePipes } from '@nestjs/common';
import { ZodValidationPipe } from '@ecommerce-backend/shared';
import { AppService } from './app.service';
import {
  type CreateUserDto,
  createUserSchema,
} from '../common/dtos/create-user.dto';

@Controller('')
export class AppController {
  constructor(private readonly userService: AppService) {}

  @Post('register')
  @UsePipes(new ZodValidationPipe(createUserSchema))
  async register(@Body() body: CreateUserDto) {
    return this.userService.register(body);
  }
}
