<script lang="ts">
  import { page } from '$app/state';

  const links = [
    { href: '/shop', label: 'shop' },
    { href: '/about', label: 'about' },
  ];
</script>

<nav class="sticky top-0 z-50 bg-sand/90 backdrop-blur-sm border-b border-border">
  <div class="mx-auto max-w-6xl flex items-center justify-between px-4 py-3">
    <a href="/" class="font-display text-xl font-semibold text-ink tracking-tight">
      retrosnack
    </a>

    <div class="flex items-center gap-6">
      {#each links as link}
        <a
          href={link.href}
          class="text-sm transition-colors {page.url.pathname.startsWith(link.href)
            ? 'text-ink font-medium'
            : 'text-ink-muted hover:text-ink'}"
        >
          {link.label}
        </a>
      {/each}
    </div>
  </div>
</nav>
