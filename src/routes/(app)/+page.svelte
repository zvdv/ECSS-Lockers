<script lang="ts">
  import Button from "$lib/components/Button.svelte";
  import Flash from "$lib/components/Flash.svelte";
  import Input from "$lib/components/Input.svelte";
  import Logo from "$lib/components/Logo.svelte";
  import Title from "$lib/components/Title.svelte";
  import { superForm } from "$lib/form.client";

  export let data;
  const { form, delayed, enhance, errors, constraints, message } = superForm(
    data.form
  );

  const logoSize = "6rem";
</script>

<Title />
<main style:--logo-size={logoSize}>
  <Logo size={logoSize} />
  <h1>Locker Registration</h1>
  {#if $message}
    <Flash {...$message} />
  {/if}
  <form method="post" use:enhance>
    <Input
      type="text"
      label="Name or Club"
      id="name"
      name="name"
      bind:value={$form.name}
      errors={$errors.name}
      {...$constraints.name}
    />
    <Input
      type="email"
      label="UVic Email"
      id="email"
      name="email"
      bind:value={$form.email}
      errors={$errors.email}
      {...$constraints.email}
    />
    <Input
      type="text"
      label="Locker Number"
      id="locker"
      name="locker"
      list="available-lockers"
      bind:value={$form.locker}
      errors={$errors.locker}
      {...$constraints.locker}
    />
    <datalist id="available-lockers">
      {#each data.availableLockers as locker}
        <option value={locker} />
      {/each}
    </datalist>
    <div class="w-full text-sm font-medium text-neutral-800">
      <input class="cursor-pointer" type="checkbox" id="agree" required />
      <label class="cursor-pointer select-none" for="agree">
        I have read and agree to the
        <a
          class="text-neutral-600 underline"
          href="https://onlineacademiccommunity.uvic.ca/ess/lockers/"
        >
          terms of use
        </a>
      </label>
    </div>
    <Button loading={$delayed}>Register</Button>
  </form>
  <footer>
    Already have an account?
    <a href="./login">Log in.</a>
  </footer>
</main>

<style lang="postcss">
  main {
    position: relative;
    top: calc((var(--logo-size) + var(--main-gap)) / -3);
  }
</style>
