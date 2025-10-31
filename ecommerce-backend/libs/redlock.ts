import {
  Inject,
  Injectable,
  OnModuleDestroy,
  OnModuleInit,
} from '@nestjs/common';
import { Redis } from 'ioredis';
import Redlock from 'redlock';

@Injectable()
export class RedlockService implements OnModuleInit, OnModuleDestroy {
  private redlock: Redlock;

  constructor(@Inject('REDIS_CLIENT') private readonly redis: Redis) {}

  onModuleInit() {
    this.redlock = new Redlock([this.redis], {
      retryCount: 3,
      retryDelay: 200,
      retryJitter: 100,
    });

    this.redlock.on('error', (err: unknown) => {
      // This fires when a lock cannot be acquired or extended
      console.error('🔒 Redlock error:', err);
    });

    console.log('✅ Redlock initialized');
  }

  async acquireLock(resource: string, ttl = 3000) {
    // resource: a unique key, like `locks:product:${variantId}`
    try {
      const lock = await this.redlock.acquire([resource], ttl);
      return lock;
    } catch (err) {
      throw new Error(`Failed to acquire lock for ${resource}`);
    }
  }

  async releaseLock(lock: Redlock.Lock) {
    try {
      await lock.release();
    } catch (err) {
      console.error('Error releasing lock:', err);
    }
  }

  onModuleDestroy() {
    // clean up connections if needed
  }
}
