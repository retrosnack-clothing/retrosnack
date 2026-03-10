<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';

	onMount(async () => {
		if ('serviceWorker' in navigator) {
			const { registerSW } = await import('virtual:pwa-register');
			registerSW({ immediate: true });
		}
	});
</script>

<slot />
