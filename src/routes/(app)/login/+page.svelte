<script lang="ts">
  import Button from "$lib/components/Button.svelte";
  import Flash from "$lib/components/Flash.svelte";
  import Input from "$lib/components/Input.svelte";
  import Logo from "$lib/components/Logo.svelte";
  import Title from "$lib/components/Title.svelte";
  import { superForm } from "$lib/form.client";

  export let data;
  const { form, delayed, enhance, errors, message } = superForm(data.form);

  const logoSize = "6rem";
</script>

<Title value="Login" />
<main style:--logo-size={logoSize}>
  <Logo size="6rem" />
  <h1>Locker Registration</h1>
  {#if $message}
    <Flash {...$message} />
  {/if}
  <form method="post" use:enhance>
    <!-- type is text to allow secret admin login -->
    <Input
      type="text"
      label="Email"
      id="email"
      name="email"
      bind:value={$form.email}
      errors={$errors.email}
    />
    <Button loading={$delayed}>Login</Button>
  </form>
  <footer>
    Don't have a locker?
    <a href="./">Register</a>
  </footer>
</main>

<style lang="postcss">
  main {
    position: relative;
    top: calc((var(--logo-size) + var(--main-gap)) / -3);
  }
</style>
