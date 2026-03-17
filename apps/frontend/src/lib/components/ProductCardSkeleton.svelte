<script lang="ts">
    import Skeleton from './Skeleton.svelte';
</script>

<div class="rounded-xl bg-sand-light overflow-hidden" style="box-shadow: var(--shadow-soft)">
    <Skeleton class="aspect-[3/4] w-full rounded-none" />
    <div class="px-3.5 py-3.5 space-y-2">
        <Skeleton class="h-4 w-3/4" />
        <div class="flex items-center justify-between">
            <Skeleton class="h-3 w-1/3" />
            <Skeleton class="h-3.5 w-1/5" />
        </div>
    </div>
</div>
