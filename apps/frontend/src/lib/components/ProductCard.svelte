<script lang="ts">
  interface Props {
    id: string;
    title: string;
    price: number;
    size: string;
    image: string;
    sold?: boolean;
  }

  let { id, title, price, size, image, sold = false }: Props = $props();
</script>

<a href="/shop/{id}" class="group block {sold ? 'opacity-60' : ''}">
  <div class="relative aspect-[3/4] overflow-hidden rounded-lg bg-sand-dark">
    <img
      src={image}
      alt={title}
      class="h-full w-full object-cover transition-transform duration-300 {sold ? '' : 'group-hover:scale-105'}"
    />
    {#if sold}
      <div class="absolute inset-0 flex items-center justify-center">
        <span class="bg-ink/80 text-sand text-xs font-medium px-4 py-1.5 rounded-full">
          sold
        </span>
      </div>
    {/if}
  </div>

  <div class="mt-3 space-y-1">
    <h3 class="text-sm font-medium truncate">{title}</h3>
    <div class="flex items-center justify-between text-sm">
      <span class="text-ink-muted">{size}</span>
      <span class="font-semibold">${(price / 100).toFixed(2)}</span>
    </div>
  </div>
</a>
