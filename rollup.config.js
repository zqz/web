// rollup.config.js
import * as fs from 'fs';
import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import autoPreprocess from 'svelte-preprocess';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';

const production = !process.env.ROLLUP_WATCH;

export default {
  watch: {
    chokidar: false,
  },
  input: 'src/main.js',
  output: {
    sourcemap: !production,
    name: 'app',
    file: 'public/build/bundle.js',
    format: 'iife'
  },
  plugins: [
    svelte({
      include: 'src/**/*.svelte',
      hydratable: true,
      preprocess: autoPreprocess(),
      css: function (css) {
        css.write('public/build/main.css');
      },

      onwarn: (warning, handler) => {
        handler(warning);
      }
    }),
    resolve({
      browser: true,
      dedupe: ['svelte']
    }),
    commonjs(),
    !production && livereload({watch: 'public'}),
    production && terser(),
    !production && serve()
  ]
};

function serve() {
  let started = false;

  return {
    writeBundle() {
      if (!started) {
        started = true;

        require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
          stdio: ['ignore', 'inherit', 'inherit'],
          shell: true
        });
      }
    }
  };
}
