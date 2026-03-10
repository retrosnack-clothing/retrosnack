<script lang="ts">
  import ProductCard from '$lib/components/ProductCard.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { api } from '$lib/api';
  import type { Product } from '$lib/api';

  let products = $state<Product[]>([]);
  let loading = $state(true);
  let error = $state('');

  async function loadProducts() {
    try {
      products = await api.products.list();
    } catch (e) {
      error = 'failed to load products';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadProducts();
  });
</script>

<svelte:head>
  <title>shop — retrosnack clothing</title>
</svelte:head>

<section class="mx-auto max-w-6xl px-4 py-12">
  <h1 class="text-2xl md:text-3xl font-semibold mb-8">the rack</h1>

  {#if loading}
    <p class="text-center text-ink-muted py-16">loading...</p>
  {:else if error}
    <p class="text-center text-ink-muted py-16">{error}</p>
  {:else if products.length === 0}
    <p class="text-center text-ink-muted py-16">no items on the rack right now — check back soon.</p>
  {:else}
    <ProductGrid>
      {#each products as product (product.id)}
        <ProductCard {product} />
      {/each}
    </ProductGrid>
  {/if}
</section>
