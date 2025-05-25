import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [svelte()],

  resolve: {
    alias: {
      '$lib': path.resolve(__dirname, './src/lib'),
      '$utils': path.resolve(__dirname, './src/utils')
    }
  },

  build: {
    lib: {
      entry: 'src/widget.ts',
      name: 'Widget',
      fileName: (format) => `widget.${format}.js`
    },
  },

  server: {
    cors: true,
    host: true,
    port: 5173,
    headers: {
      'Access-Control-Allow-Origin': '*'
    }
  }
});
