<script lang="ts">
  import Back from "$lib/components/Back.svelte";
  import Button from "$lib/components/Button.svelte";
</script>

<main>
  <h1>Deregister locker</h1>
  <p class="text-center">
    This locker will no longer be yours. Please remove any locks and belongings
    from the locker.
  </p>
  <form method="post">
    <Button>Deregister</Button>
  </form>
</main>
<Back />
