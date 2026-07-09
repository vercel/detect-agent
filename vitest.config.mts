import { defineConfig } from 'vite';

export default defineConfig({
  test: {
    globals: true,
    // Use of process.chdir prohibits usage of the default "threads". https://vitest.dev/config/#forks
    pool: 'forks',
    include: ['src/**/*.{test,spec}.?(c|m)[jt]s?(x)'],
    env: {
      // Vitest suppresses color output when `process.env.CI` is true.
      // https://github.com/vitest-dev/vitest/issues/2732
      FORCE_COLOR: '1',
    },
  },
});
