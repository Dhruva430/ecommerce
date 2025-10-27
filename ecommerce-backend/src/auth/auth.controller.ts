import { Post } from '@nestjs/common';
import { Role } from '@prisma/client';
import { Roles } from 'src/common/roles/roles.decorator';

export class Authcontroller {
  @Post()
  @Roles(Role.ADMIN)
  create() {
    return 'this is a test';
  }
}
