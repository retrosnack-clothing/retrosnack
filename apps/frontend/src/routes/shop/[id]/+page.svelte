<script lang="ts">
  import { page } from '$app/state';
  import { mockProducts } from '$lib/mock-data';
  import { cart } from '$lib/stores/cart.svelte';

  const product = $derived(mockProducts.find((p) => p.id === page.params.id));
  const inCart = $derived(cart.items.some((i) => i.id === page.params.id));

  function addToCart() {
    if (!product || inCart) return;
    cart.add({
      id: product.id,
      title: product.title,
      price: product.price,
      size: product.size,
      image: product.image,
    });
  }
</script>

<svelte:head>
  <title>{product ? product.title : 'not found'} — retrosnack clothing</title>
</svelte:head>

{#if product}
  <article class="mx-auto max-w-4xl px-4 py-12">
    <div class="grid md:grid-cols-2 gap-8 md:gap-12">
      <div class="aspect-[3/4] overflow-hidden rounded-lg bg-sand-dark">
        <img
          src={product.image}
          alt={product.title}
          class="h-full w-full object-cover"
        />
      </div>

      <div class="flex flex-col justify-center gap-6">
        <div>
          <p class="text-sm text-ink-muted mb-1">{product.brand}</p>
          <h1 class="text-2xl md:text-3xl font-semibold">{product.title}</h1>
        </div>

        <div class="flex items-center gap-3 text-sm text-ink-muted">
          <span class="bg-tag px-3 py-1 rounded-full">size {product.size}</span>
          <span class="bg-tag px-3 py-1 rounded-full">{product.condition}</span>
        </div>

        <p class="text-3xl font-semibold">
          ${(product.price / 100).toFixed(2)}
        </p>

        <div class="flex flex-col gap-3">
          {#if product.sold}
            <button
              disabled
              class="bg-ink/40 text-sand px-6 py-3 rounded-full text-sm font-medium cursor-not-allowed"
            >
              sold
            </button>
          {:else if inCart}
            <a
              href="/cart"
              class="bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium text-center hover:bg-ink/85 transition-colors"
            >
              in your bag — view bag
            </a>
          {:else}
            <button
              onclick={addToCart}
              class="bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
            >
              add to bag
            </button>
          {/if}

          <a
            href="https://instagram.com/retrosnack.shop"
            target="_blank"
            rel="noopener noreferrer"
            class="border border-border px-6 py-3 rounded-full text-sm font-medium text-center hover:bg-sand-dark transition-colors"
          >
            message on instagram
          </a>
        </div>
      </div>
    </div>
  </article>
{:else}
  <div class="mx-auto max-w-4xl px-4 py-24 text-center">
    <h1 class="text-2xl font-semibold mb-2">item not found</h1>
    <p class="text-ink-muted mb-6">this piece may have already found a new home.</p>
    <a href="/shop" class="text-accent hover:text-accent-hover transition-colors">
      back to the rack &rarr;
    </a>
  </div>
{/if}
