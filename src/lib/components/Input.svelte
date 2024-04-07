<script lang="ts">
  import type { HTMLInputTypeAttribute } from "svelte/elements";

  export let id: string;
  export let label: string;
  export let name: string;
  export let value = "";
  export let errors: string[] = [];
  // exclude number types, we can make another component in the future
  export let type: Exclude<HTMLInputTypeAttribute, "number" | "range"> = "text";
  /** Makes the element appear disabled, although it's still submitted */
  export let disabled = false;

  // Svelte doesn't allow dynamic type binding, so we have to force it
  // https://stackoverflow.com/a/57393751
  function handleInput(
    e: Event & {
      currentTarget: EventTarget & HTMLInputElement;
    }
  ) {
    value = e.currentTarget.value;
    // if (type.match(/^(number|range)$/)) {
    //   value = +e.currentTarget.value;
    // } else {
    //   value = e.currentTarget.value;
    // }
  }
</script>

<div class="container">
  <label for={id}>
    {label}
  </label>
  <input
    class:error={(errors ?? []).length}
    {type}
    {id}
    {name}
    {disabled}
    required
    {value}
    on:input={handleInput}
    {...$$restProps}
  />
  {#if disabled}
    <input class="hidden" bind:value {name} />
  {/if}
  {#each errors ?? [] as error}
    <span class="error-msg">* {error}</span>
  {/each}
</div>

<style lang="postcss">
  .container {
    width: 100%;
    max-width: theme(maxWidth.xs);
  }
  label {
    color: theme(colors.neutral.600);
    font-size: theme(fontSize.base);
    font-weight: theme(fontWeight.medium);
  }
  input {
    display: block;
    height: theme(space.10);
    width: 100%;
    margin: auto;
    margin-top: theme(space.1);
    padding: theme(padding[0.5]) theme(padding.2);
    box-sizing: border-box;
    border: 1px solid theme(colors.neutral.300);
    @apply rounded;
    background: theme(colors.neutral.50);
    color: theme(colors.neutral.700);
    font-size: theme(fontSize.base);
  }
  input.error {
    border-color: theme(colors.red.300);
  }
  input.hidden {
    display: none;
  }
  input:disabled {
    background: theme(colors.neutral.200);
    color: theme(colors.neutral.500);
  }
  .error-msg {
    color: theme(colors.red.600);
    font-size: theme(fontSize.sm);
    font-weight: theme(fontWeight.medium);
  }
</style>
