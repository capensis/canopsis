import { fileURLToPath } from 'node:url';

import { mergeConfig, defineConfig, configDefaults } from 'vitest/config';

import viteConfig from './vite.config.js';

const TEST_EXCLUDED_PLUGINS = new Set([
  'vite-plugin-inspect',
  'vite-plugin-vue-devtools',
  'vite:basic-ssl',
]);
const STYLE_RE = /\.(?:css|less|sass|scss|styl|stylus|pcss|postcss)(?:$|\?)/;
const ASSET_RE = /\.(?:avif|gif|jpe?g|png|svg|webp)(?:$|\?)/;

const testAssetsPlugin = {
  name: 'canopsis:test-assets',
  enforce: 'pre',
  load(id) {
    if (id.includes('?vue&type=style') || STYLE_RE.test(id)) {
      return '';
    }

    if (ASSET_RE.test(id)) {
      const fileName = id.split('?')[0].split('/').at(-1);

      return `export default ${JSON.stringify(fileName)};`;
    }

    return null;
  },
};

const testViteConfig = {
  ...viteConfig,
  plugins: [
    testAssetsPlugin,
    ...viteConfig.plugins
      .flat()
      .filter(plugin => plugin && !TEST_EXCLUDED_PLUGINS.has(plugin.name)),
  ],
};

export default mergeConfig(
  testViteConfig,
  defineConfig({
    test: {
      environment: 'jsdom',
      exclude: [...configDefaults.exclude, 'e2e/**'],
      globalSetup: ['./__tests__/global-setup.js'],
      globals: true,
      pool: 'threads',
      root: fileURLToPath(new URL('./', import.meta.url)),
      setupFiles: ['./__tests__/setup.js'],
      silent: 'passed-only',
      deps: {
        web: {
          transformAssets: false,
          transformCss: false,
        },
      },
      server: {
        deps: {
          inline: ['vuetify'],
        },
      },
    },
  }),
);
