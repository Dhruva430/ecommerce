import { Body, Controller, Post, UsePipes } from '@nestjs/common';
import { z } from 'zod';
import { ZodValidationPipe } from '@ecommerce-backend/shared';
import { AppService } from './app.service';

const createUserSchema = z.object({
  name: z.string().min(2),
  email: z.string().email(),
  password: z.string().min(6),
});

@Controller('users')
export class UserController {
  constructor(private readonly userService: AppService) {}

  @Post()
  @UsePipes(new ZodValidationPipe(createUserSchema))
  async createUser(@Body() dto: z.infer<typeof createUserSchema>) {
    return this.userService.createUser(dto);
  }
}
