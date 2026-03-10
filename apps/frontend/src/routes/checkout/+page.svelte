<script lang="ts">
  import { cart } from '$lib/stores/cart.svelte';
  import { api } from '$lib/api';

  let submitting = $state(false);
  let error = $state('');

  async function handleCheckout() {
    if (submitting || cart.items.length === 0) return;

    submitting = true;
    error = '';

    try {
      // create order from cart items
      const order = await api.orders.create(
        cart.items.map((item) => ({
          variant_id: item.variantId,
          quantity: 1,
          price_cents: item.priceCents,
        }))
      );

      // create square checkout session and redirect
      const session = await api.checkout.create(order.id);
      cart.clear();
      window.location.href = session.url;
    } catch {
      error = 'something went wrong — please try again.';
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>checkout — retrosnack clothing</title>
</svelte:head>

<section class="mx-auto max-w-2xl px-4 py-12">
  <h1 class="text-2xl md:text-3xl font-semibold mb-8">checkout</h1>

  {#if cart.items.length === 0}
    <div class="text-center py-16">
      <p class="text-ink-muted mb-6">your bag is empty.</p>
      <a
        href="/shop"
        class="inline-block bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
      >
        browse the rack
      </a>
    </div>
  {:else}
    <div class="bg-sand-light rounded-lg p-6 border border-border mb-8">
      <h2 class="text-lg font-medium mb-4">order summary</h2>

      <div class="space-y-4 mb-6">
        {#each cart.items as item (item.variantId)}
          <div class="flex gap-3">
            {#if item.image}
              <img
                src={item.image}
                alt={item.title}
                class="w-14 h-18 object-cover rounded bg-sand-dark"
              />
            {:else}
              <div class="w-14 h-18 rounded bg-sand-dark"></div>
            {/if}
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">{item.title}</p>
              <p class="text-xs text-ink-muted">{item.size}{item.color ? ` · ${item.color}` : ''}</p>
            </div>
            <p class="text-sm font-semibold shrink-0">${(item.priceCents / 100).toFixed(2)}</p>
          </div>
        {/each}
      </div>

      <div class="border-t border-border pt-4 space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-ink-muted">subtotal</span>
          <span>${(cart.totalCents / 100).toFixed(2)}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-ink-muted">shipping</span>
          <span class="text-ink-muted">calculated at payment</span>
        </div>
        <div class="flex justify-between font-semibold pt-2 border-t border-border">
          <span>total</span>
          <span>${(cart.totalCents / 100).toFixed(2)}</span>
        </div>
      </div>
    </div>

    {#if error}
      <p class="text-accent text-sm mb-4">{error}</p>
    {/if}

    <button
      onclick={handleCheckout}
      disabled={submitting}
      class="w-full bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium transition-colors {submitting ? 'opacity-60 cursor-not-allowed' : 'hover:bg-ink/85'}"
    >
      {submitting ? 'redirecting to payment...' : 'pay with square'}
    </button>

    <p class="text-xs text-ink-muted text-center mt-3">
      you'll enter shipping and payment details on the next page
    </p>
  {/if}
</section>
