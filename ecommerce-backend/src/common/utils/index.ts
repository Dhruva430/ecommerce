import bycrpyt from 'bcrypt';
import { boolean } from 'zod';

export default function hashPassword(password: string) {
  const saltRounds = 10;
  const hash = bycrpyt.hashSync(password, saltRounds);
  return hash;
}
export function comparePassword(
  plainPassword: string,
  hashedPassword: string,
): boolean {
  const isMatch = bycrpyt.compareSync(plainPassword, hashedPassword);
  if (!isMatch) {
    return false;
  }
  return true;
}
