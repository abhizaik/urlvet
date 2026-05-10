import adapterAuto from '@sveltejs/adapter-auto';
import adapterStatic from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),

  kit: {
    // Docker builds set DOCKER_BUILD=1 → adapter-static (served by Caddy).
    // All other environments (Vercel, local dev) → adapter-auto.
    adapter:
      process.env.DOCKER_BUILD === '1' ? adapterStatic({ fallback: '200.html' }) : adapterAuto(),
  },
};

export default config;
