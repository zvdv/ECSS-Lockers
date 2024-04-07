<script lang="ts">
  import Badge from "./Badge.svelte";

  const item = ["available", "claimed", "expired"] as const;
  export let hidden = false;
  export let filter = Object.fromEntries(item.map((status) => [status, true]));
</script>

<div class="absolute left-0 right-0 text-left" class:hidden>
  <div
    class="relative left-[-1px] top-1 w-[calc(100%+2px)] border border-black bg-neutral-200 p-1"
  >
    {#each item as key}
      <div>
        <input type="checkbox" bind:checked={filter[key]} value={key} />
        <span
          on:click={() => {
            filter[key] = !filter[key];
          }}
          on:keypress={() => {
            filter[key] = !filter[key];
          }}
          class="cursor-pointer select-none font-normal"
        >
          <Badge status={key}>{key}</Badge>
        </span>
      </div>
    {/each}
  </div>
</div>
