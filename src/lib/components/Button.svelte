<script lang="ts">
  import Loading from "./Loading.svelte";

  export let compact = false;
  export let loading = false;
</script>

<button class:compact class:loading on:click>
  <div class="spinner">
    <Loading size="1.5rem" />
  </div>
  <slot />
</button>

<style lang="postcss">
  button {
    padding: theme(space.3) theme(space.6);

    position: relative;
    display: inline-flex;
    align-items: center;
    justify-content: center;

    border-color: transparent;
    border-radius: theme(borderRadius.lg);
    color: theme(colors.neutral.50);
    background-color: theme(colors.ess.500);

    font-weight: theme(fontWeight.bold);
    font-size: theme(fontSize.sm);

    transition: background-color 0.2s ease-in-out;
  }
  button:hover {
    background-color: theme(colors.ess.400);
  }
  .compact {
    padding: theme(space.2) theme(space.4);
  }
  .loading {
    color: transparent;
  }
  div {
    display: none;
  }
  .loading > div {
    position: absolute;
    inset: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
</style>
