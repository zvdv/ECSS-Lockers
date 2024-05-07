<script lang="ts">
  import { page } from "$app/stores";
  import Shell from "$lib/components/Shell.svelte";
</script>

<Shell>
  <main class="max-w-sm">
    <h1 class="text-3xl font-bold">
      <span class="underline decoration-red-500 decoration-wavy"
        >{$page.status}</span
      >
      - {$page.error?.message}
    </h1>
    <p>
      Something went wrong. If you continue to see this error, please contact
      the ECSS Director of IT.
    </p>
  </main>
</Shell>
