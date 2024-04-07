<script lang="ts">
  import { enhance } from "$app/forms";
  import Button from "$lib/components/Button.svelte";
  import Menu from "./Menu.svelte";

  export let data;

  let logoutLoading = false;
</script>

<main>
  <h1>Registered Lockers</h1>
  <Menu rows={data.locker} />
  <form
    method="post"
    action="?/logout"
    use:enhance={() => {
      logoutLoading = true;
      return ({ update }) => {
        logoutLoading = false;
        update();
      };
    }}
  >
    <Button loading={logoutLoading}>Logout</Button>
  </form>
</main>
