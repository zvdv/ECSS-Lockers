<script lang="ts">
  export let id: string;
  export let name: string;
  export let files: FileList;
  export let errors: string[] = [];
</script>

<div class="w-full">
  <div
    class="relative flex w-full flex-col items-center justify-center gap-2 break-words rounded border border-neutral-300 bg-neutral-50 p-10"
    class:error={errors?.length}
  >
    {#if files?.length > 0}
      <h1 class="text-center text-lg font-bold">
        {files[0].name}
      </h1>
    {:else}
      <h1 class="text-center text-xl font-bold text-neutral-500">
        Drop file here
      </h1>
      <span class="font-bold text-neutral-500">or</span>
    {/if}
    <label
      class="cursor-pointer rounded bg-ess-500 px-3 py-1.5 text-sm font-medium text-white
    "
      for={id}>Click to open file browser</label
    >
    <input
      class="absolute inset-0 cursor-pointer opacity-0"
      type="file"
      {id}
      {name}
      accept=".csv,.xlsx"
      bind:files
      {...$$props}
    />
  </div>
  {#each errors ?? [] as error}
    <span class="text-sm font-medium text-red-600">* {error}</span>
  {/each}
</div>

<style lang="postcss">
  .error {
    @apply border-red-300;
  }
</style>
