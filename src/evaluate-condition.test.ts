import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import mockFs from 'mock-fs';
import { evaluateCondition } from './evaluate-condition';
import type { Condition } from './types';

describe('evaluateCondition', () => {
  beforeEach(() => {
    vi.unstubAllEnvs();
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    mockFs.restore();
  });

  describe('env_set', () => {
    it('is true when the variable has a non-empty value', async () => {
      vi.stubEnv('SOME_VAR', '1');
      expect(
        await evaluateCondition({ type: 'env_set', name: 'SOME_VAR' })
      ).toBe(true);
    });

    it('is false when the variable is unset', async () => {
      vi.stubEnv('SOME_VAR', '');
      expect(
        await evaluateCondition({ type: 'env_set', name: 'SOME_VAR' })
      ).toBe(false);
    });

    it('is false when the variable is set to an empty string', async () => {
      // `env_set` uses Boolean(), so an empty string is falsy.
      vi.stubEnv('SOME_VAR', '');
      expect(
        await evaluateCondition({ type: 'env_set', name: 'SOME_VAR' })
      ).toBe(false);
    });
  });

  describe('env_value', () => {
    it('is true when the value matches exactly', async () => {
      vi.stubEnv('ROLE', 'agent-exec');
      expect(
        await evaluateCondition({
          type: 'env_value',
          name: 'ROLE',
          value: 'agent-exec',
        })
      ).toBe(true);
    });

    it('is false when the value differs', async () => {
      vi.stubEnv('ROLE', 'something-else');
      expect(
        await evaluateCondition({
          type: 'env_value',
          name: 'ROLE',
          value: 'agent-exec',
        })
      ).toBe(false);
    });

    it('is false when the variable is unset', async () => {
      expect(
        await evaluateCondition({
          type: 'env_value',
          name: 'ROLE',
          value: 'agent-exec',
        })
      ).toBe(false);
    });

    it('does not treat an empty target value as matching an unset variable', async () => {
      // process.env access for an unset var is undefined, not '', so an
      // `env_value` looking for '' must not match.
      expect(
        await evaluateCondition({ type: 'env_value', name: 'ROLE', value: '' })
      ).toBe(false);
    });
  });

  describe('file_exists', () => {
    it('is true when the path exists', async () => {
      mockFs({ '/opt/.devin': mockFs.directory({ mode: 0o755 }) });
      expect(
        await evaluateCondition({ type: 'file_exists', path: '/opt/.devin' })
      ).toBe(true);
    });

    it('is false when the path does not exist', async () => {
      mockFs({});
      expect(
        await evaluateCondition({ type: 'file_exists', path: '/opt/.devin' })
      ).toBe(false);
    });

    it('is false when access throws (e.g. unreadable parent)', async () => {
      mockFs({ '/opt': mockFs.directory({ mode: 0o000 }) });
      expect(
        await evaluateCondition({ type: 'file_exists', path: '/opt/.devin' })
      ).toBe(false);
    });
  });

  describe('env_matches', () => {
    it('is true when the value matches the pattern', async () => {
      vi.stubEnv('PATH', '/usr/bin:/home/me/.pi/agent/bin');
      expect(
        await evaluateCondition({
          type: 'env_matches',
          name: 'PATH',
          pattern: '\\.pi[\\\\/]agent',
        })
      ).toBe(true);
    });

    it('matches a backslash path separator too', async () => {
      vi.stubEnv('PATH', 'C:\\Users\\me\\.pi\\agent\\bin');
      expect(
        await evaluateCondition({
          type: 'env_matches',
          name: 'PATH',
          pattern: '\\.pi[\\\\/]agent',
        })
      ).toBe(true);
    });

    it('is false when the value does not match', async () => {
      vi.stubEnv('PATH', '/usr/bin:/usr/local/bin');
      expect(
        await evaluateCondition({
          type: 'env_matches',
          name: 'PATH',
          pattern: '\\.pi[\\\\/]agent',
        })
      ).toBe(false);
    });

    it('is false when the variable is unset', async () => {
      vi.stubEnv('TERM_PROGRAM', '');
      expect(
        await evaluateCondition({
          type: 'env_matches',
          name: 'TERM_PROGRAM',
          pattern: 'kiro',
        })
      ).toBe(false);
    });

    it('is false when the pattern is invalid rather than throwing', async () => {
      vi.stubEnv('TERM_PROGRAM', 'kiro');
      expect(
        await evaluateCondition({
          type: 'env_matches',
          name: 'TERM_PROGRAM',
          pattern: '(',
        })
      ).toBe(false);
    });
  });

  describe('no_tty', () => {
    const original = process.stdout.isTTY;
    afterEach(() => {
      process.stdout.isTTY = original;
    });

    it('is true when stdout is not a TTY', async () => {
      process.stdout.isTTY = false;
      expect(await evaluateCondition({ type: 'no_tty' })).toBe(true);
    });

    it('is false when stdout is a TTY', async () => {
      process.stdout.isTTY = true;
      expect(await evaluateCondition({ type: 'no_tty' })).toBe(false);
    });
  });

  describe('anyOf', () => {
    it('is true when at least one child is true', async () => {
      vi.stubEnv('B', '1');
      expect(
        await evaluateCondition({
          type: 'anyOf',
          conditions: [
            { type: 'env_set', name: 'A' },
            { type: 'env_set', name: 'B' },
          ],
        })
      ).toBe(true);
    });

    it('is false when all children are false', async () => {
      expect(
        await evaluateCondition({
          type: 'anyOf',
          conditions: [
            { type: 'env_set', name: 'A' },
            { type: 'env_set', name: 'B' },
          ],
        })
      ).toBe(false);
    });

    it('is false for an empty condition list', async () => {
      expect(await evaluateCondition({ type: 'anyOf', conditions: [] })).toBe(
        false
      );
    });

    it('is true when a later child is satisfied but an earlier one is not', async () => {
      vi.stubEnv('A', '');
      vi.stubEnv('B', '1');
      expect(
        await evaluateCondition({
          type: 'anyOf',
          conditions: [
            { type: 'env_set', name: 'A' },
            { type: 'env_set', name: 'B' },
          ],
        })
      ).toBe(true);
    });
  });

  describe('allOf', () => {
    it('is true when all children are true', async () => {
      vi.stubEnv('A', '1');
      vi.stubEnv('B', '1');
      expect(
        await evaluateCondition({
          type: 'allOf',
          conditions: [
            { type: 'env_set', name: 'A' },
            { type: 'env_set', name: 'B' },
          ],
        })
      ).toBe(true);
    });

    it('is false when any child is false', async () => {
      vi.stubEnv('A', '1');
      expect(
        await evaluateCondition({
          type: 'allOf',
          conditions: [
            { type: 'env_set', name: 'A' },
            { type: 'env_set', name: 'B' },
          ],
        })
      ).toBe(false);
    });

    it('is true for an empty condition list (vacuous truth)', async () => {
      expect(await evaluateCondition({ type: 'allOf', conditions: [] })).toBe(
        true
      );
    });

    it('is false when a later child is unsatisfied but an earlier one is not', async () => {
      vi.stubEnv('A', '1');
      vi.stubEnv('B', '');
      expect(
        await evaluateCondition({
          type: 'allOf',
          conditions: [
            { type: 'env_set', name: 'A' },
            { type: 'env_set', name: 'B' },
          ],
        })
      ).toBe(false);
    });
  });

  describe('nesting', () => {
    it('evaluates the cowork-style allOf(env_set, anyOf(...)) tree', async () => {
      vi.stubEnv('CLAUDE_CODE_IS_COWORK', '1');
      vi.stubEnv('CLAUDECODE', '1');
      const condition: Condition = {
        type: 'allOf',
        conditions: [
          { type: 'env_set', name: 'CLAUDE_CODE_IS_COWORK' },
          {
            type: 'anyOf',
            conditions: [
              { type: 'env_set', name: 'CLAUDECODE' },
              { type: 'env_set', name: 'CLAUDE_CODE' },
            ],
          },
        ],
      };
      expect(await evaluateCondition(condition)).toBe(true);
    });

    it('is false when the nested anyOf has no satisfied child', async () => {
      vi.stubEnv('CLAUDE_CODE_IS_COWORK', '1');
      vi.stubEnv('CLAUDECODE', '');
      vi.stubEnv('CLAUDE_CODE', '');
      const condition: Condition = {
        type: 'allOf',
        conditions: [
          { type: 'env_set', name: 'CLAUDE_CODE_IS_COWORK' },
          {
            type: 'anyOf',
            conditions: [
              { type: 'env_set', name: 'CLAUDECODE' },
              { type: 'env_set', name: 'CLAUDE_CODE' },
            ],
          },
        ],
      };
      expect(await evaluateCondition(condition)).toBe(false);
    });
  });

  describe('unknown condition type', () => {
    it('returns false for an unrecognized type', async () => {
      // Simulates a condition type not covered by the switch — the `default`
      // branch must fail closed rather than throw.
      const unknown = { type: 'not_a_real_type' } as unknown as Condition;
      expect(await evaluateCondition(unknown)).toBe(false);
    });
  });
});
