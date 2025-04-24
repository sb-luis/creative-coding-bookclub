// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';
import { viteStaticCopy } from 'vite-plugin-static-copy'

// https://astro.build/config
export default defineConfig({
  vite: {
    plugins: [
      tailwindcss(),
      viteStaticCopy({
        targets: [
          {
            src: './src/members/**/*.[jt]s',
            dest: 'sketch-files'
          }
        ]
      })
    ],
  },
  devToolbar: {
    enabled: false,
  },
});
