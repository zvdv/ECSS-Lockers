<script lang="ts">
  import { enhance } from "$app/forms";
  import Back from "$lib/components/Back.svelte";
  import Button from "$lib/components/Button.svelte";

  export let data;
  let renewLoading = false;
</script>

<main>
  <h1>Locker {data.locker}</h1>
  <div>
    <p><b>Register for</b>: {data.name}</p>
    <p>
      <b>Status</b>:
      {#if data.status === "claimed"}
        <span class="status green">Claimed</span>
      {:else if data.status === "expired"}
        <span class="status red">Expired</span>
      {/if}
    </p>
  </div>
  <div class="buttons">
    {#if data.status === "expired"}
      <form
        action="?/renew"
        method="post"
        use:enhance={() => {
          renewLoading = true;
          return ({ update }) => {
            renewLoading = false;
            update();
          };
        }}
      >
        <Button loading={renewLoading}>Renew registration</Button>
      </form>
    {/if}
    <a href="{data.locker}/transfer"
      ><Button>Transfer to a different email</Button></a
    >
    <a href="{data.locker}/deregister"><Button>Deregister locker</Button></a>
  </div>
</main>
<Back />

<style lang="postcss">
  .status {
    @apply rounded p-1 text-sm font-bold;
  }
  .green {
    @apply bg-green-200 text-green-800;
  }
  .red {
    @apply bg-red-200 text-red-800;
  }
  .buttons {
    @apply flex flex-col items-center gap-2;
  }
</style>
