import { sveltekit } from '@sveltejs/kit/vite';
import { VitePWA } from 'vite-plugin-pwa';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		sveltekit(),
		VitePWA({
			registerType: 'autoUpdate',
			// manifest is served from static/manifest.json
			manifest: false,
			workbox: {
				globPatterns: ['**/*.{js,css,html,ico,png,svg,webp,woff2}'],
				runtimeCaching: [
					{
						urlPattern: /^https:\/\/.*\/api\/products/,
						handler: 'NetworkFirst',
						options: {
							cacheName: 'api-products',
							expiration: { maxEntries: 100, maxAgeSeconds: 300 }
						}
					}
				]
			}
		})
	]
});
