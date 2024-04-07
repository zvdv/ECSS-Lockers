<script lang="ts">
  export let value = "";
  const base = "ESS Locker Registration";
</script>

<svelte:head>
  {#if value}
    <title>{value} - {base}</title>
  {:else}
    <title>{base}</title>
  {/if}
</svelte:head>
