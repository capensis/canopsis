import { fileURLToPath, URL } from 'node:url';

import basicSsl from '@vitejs/plugin-basic-ssl';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';
import vuetify from 'vite-plugin-vuetify';

// https://vite.dev/config/
export default defineConfig({
  server: {
    port: 8081,
    strictPort: true,
    hmr: {
      protocol: 'wss',
      host: 'localhost',
      clientPort: 8081,
    },
  },
  plugins: [vue(), vueDevTools(), vuetify({ autoImport: true }), basicSsl()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@unit': fileURLToPath(new URL('./src/__tests__/unit', import.meta.url)),
      'cytoscape/dist/cytoscape.umd.js': 'cytoscape/dist/cytoscape.esm.mjs',
      vue: 'vue/dist/vue.esm-bundler.js',
    },
  },
  define: {
    global: 'globalThis',
  },
});
