<script lang="ts">
    import type { Product } from '$lib/api';
    import Skeleton from './Skeleton.svelte';

    interface Props {
        product: Product;
        sold?: boolean;
    }

    let { product, sold = false }: Props = $props();

    const image = product.images[0]?.url ?? '';
    let imageLoaded = $state(!image);
</script>

<a
    href="/shop/{product.id}"
    class="group block hover-lift rounded-xl bg-sand-light overflow-hidden {sold ? 'opacity-60' : ''}"
>
    <div class="relative aspect-[3/4] overflow-hidden bg-sand-dark">
        {#if image}
            {#if !imageLoaded}
                <Skeleton class="absolute inset-0 rounded-none" />
            {/if}
            <img
                src={image}
                alt={product.title}
                loading="lazy"
                width="300"
                height="400"
                onload={() => (imageLoaded = true)}
                class="h-full w-full object-cover transition-all duration-500 {sold
                    ? ''
                    : 'group-hover:scale-105'} {imageLoaded ? 'opacity-100' : 'opacity-0'}"
            />
        {/if}
        {#if sold}
            <div class="absolute inset-0 flex items-center justify-center">
                <span class="bg-ink/80 text-sand text-xs font-medium px-4 py-1.5 rounded-full">
                    sold
                </span>
            </div>
        {/if}
        <div class="absolute top-2.5 left-2.5 flex flex-col gap-1.5">
            {#if product.condition === 'new'}
                <span class="bg-ink text-sand text-xs font-medium px-2.5 py-1 rounded-full shadow-sm">
                    new
                </span>
            {/if}
            {#if product.drop}
                <span
                    class="bg-sand/90 text-ink text-xs font-medium px-2.5 py-1 rounded-full backdrop-blur-sm shadow-sm"
                >
                    {product.drop.name}
                </span>
            {/if}
        </div>
    </div>

    <div class="px-3.5 py-3.5">
        <h3 class="text-sm font-medium truncate">{product.title}</h3>
        <div class="flex items-center justify-between mt-1 text-sm">
            <span class="text-ink-muted text-xs">{product.brand}</span>
            <span class="font-semibold">${(product.price_cents / 100).toFixed(2)}</span>
        </div>
    </div>
</a>
