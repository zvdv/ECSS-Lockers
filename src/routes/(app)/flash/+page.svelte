<script lang="ts">
  import Back from "$lib/components/Back.svelte";
  import Flash from "$lib/components/Flash.svelte";

  export let data;
</script>

<main>
  <Flash {...data} />
</main>
<Back />
