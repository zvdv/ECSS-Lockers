<script lang="ts">
  export let status: "available" | "claimed" | "expired";
</script>

<span class="rounded p-1 text-sm font-bold {status}">
  <slot />
</span>

<style lang="postcss">
  .available {
    @apply bg-neutral-200 text-neutral-800;
  }
  .claimed {
    @apply bg-green-200 text-green-800;
  }
  .expired {
    @apply bg-red-200 text-red-800;
  }
</style>
