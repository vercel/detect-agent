import { access } from 'node:fs/promises';
import { constants } from 'node:fs';
import type { Condition } from './types';

/**
 * Evaluate a condition tree. `anyOf`/`allOf` are combinators over child
 * conditions; the rest are leaves. Evaluation is async because `file_exists`
 * touches the filesystem.
 */
export async function evaluateCondition(
  condition: Condition
): Promise<boolean> {
  switch (condition.type) {
    case 'env_set':
      return Boolean(process.env[condition.name]);
    case 'env_value':
      return process.env[condition.name] === condition.value;
    case 'file_exists':
      try {
        await access(condition.path, constants.F_OK);
        return true;
      } catch {
        return false;
      }
    case 'anyOf':
      for (const sub of condition.conditions) {
        if (await evaluateCondition(sub)) {
          return true;
        }
      }
      return false;
    case 'allOf':
      for (const sub of condition.conditions) {
        if (!(await evaluateCondition(sub))) {
          return false;
        }
      }
      return true;
    default:
      return false;
  }
}
