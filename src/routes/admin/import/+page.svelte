<script lang="ts">
  import Back from "$lib/components/Back.svelte";
  import Button from "$lib/components/Button.svelte";
  import Flash from "$lib/components/Flash.svelte";
  import Title from "$lib/components/Title.svelte";
  import { superForm } from "$lib/form.client";
  import FileInput from "./FileInput.svelte";

  export let data;
  const { delayed, enhance, errors, message } = superForm(data.form);

  let files: FileList;
</script>

<Title />
<main>
  <div class="text-left">
    <h2 class="text-lg font-semibold text-neutral-700">Data format</h2>
    <ul class="ml-5 list-outside list-disc">
      <li>The only columns are "locker", "name", and "email"</li>
      <li>If the locker has a name, it must also have an email</li>
    </ul>
  </div>

  <form method="post" use:enhance>
    {#if $message}
      <Flash {...$message} />
    {/if}
    <FileInput
      bind:files
      id="sheet"
      name="sheet"
      errors={$errors.sheet}
      required
    />
    <Button loading={$delayed}>Import</Button>
  </form>
</main>
<Back />

<style lang="postcss">
</style>
