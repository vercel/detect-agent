import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/index.ts'],
  format: ['cjs'],
  target: 'es2021',
  platform: 'node',
  dts: true,
  clean: true,
  sourcemap: false,
});
