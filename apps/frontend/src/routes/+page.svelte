<script lang="ts">
    import HeroSection from '$lib/components/HeroSection.svelte';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import ProductCardSkeleton from '$lib/components/ProductCardSkeleton.svelte';
    import ProductGrid from '$lib/components/ProductGrid.svelte';
    import FadeIn from '$lib/components/FadeIn.svelte';
    import FlowerDoodle from '$lib/components/FlowerDoodle.svelte';
    import { api } from '$lib/api';
    import type { Product, Drop } from '$lib/api';

    let products = $state<Product[]>([]);
    let drops = $state<Drop[]>([]);
    let loading = $state(true);

    $effect(() => {
        Promise.all([api.products.list(4, 0), api.drops.list()])
            .then(([p, d]) => {
                products = p;
                drops = d;
            })
            .catch(() => {})
            .finally(() => (loading = false));
    });

    const latestDrop = $derived(drops[0]);
</script>

<svelte:head>
    <title>retrosnack clothing — thrift finds, loved again</title>
    <meta property="og:title" content="retrosnack clothing — thrift finds, loved again" />
    <meta
        property="og:description"
        content="Curated secondhand women's clothing, accessories, and shoes."
    />
    <meta property="og:type" content="website" />
</svelte:head>

<HeroSection />

<FadeIn>
    <section class="mx-auto max-w-6xl px-4 py-16">
        <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl md:text-2xl font-semibold">latest drops</h2>
            <a href="/shop" class="text-sm text-accent hover:text-accent-hover transition-colors">
                view all &rarr;
            </a>
        </div>

        {#if loading}
            <ProductGrid>
                {#each Array(4) as _}
                    <ProductCardSkeleton />
                {/each}
            </ProductGrid>
        {:else if products.length > 0}
            <ProductGrid>
                {#each products as product (product.id)}
                    <ProductCard {product} />
                {/each}
            </ProductGrid>
        {/if}
    </section>
</FadeIn>

{#if latestDrop}
    <FadeIn>
        <section
            class="bg-sand-light"
            style="box-shadow: inset 0 1px 0 var(--color-border), inset 0 -1px 0 var(--color-border)"
        >
            <a
                href="/drops/{latestDrop.slug}"
                class="group block mx-auto max-w-6xl px-4 py-16 md:py-24 text-center relative overflow-hidden"
            >
                <FlowerDoodle
                    size={100}
                    class="absolute -right-6 -top-6 text-ink/[0.03] rotate-12 hidden md:block"
                />
                <FlowerDoodle
                    size={70}
                    class="absolute -left-4 -bottom-4 text-ink/[0.03] -rotate-12 hidden md:block"
                />
                <p class="text-xs uppercase tracking-widest text-ink-muted mb-4 relative">
                    latest collection
                </p>
                <h2 class="text-3xl md:text-5xl font-semibold relative">
                    {latestDrop.name}
                </h2>
                {#if latestDrop.description}
                    <p class="text-ink-muted mt-4 max-w-md mx-auto relative">
                        {latestDrop.description}
                    </p>
                {/if}
                <span class="btn-primary inline-block mt-6 px-6 py-3 relative">
                    shop the drop &rarr;
                </span>
            </a>
        </section>
    </FadeIn>
{/if}
