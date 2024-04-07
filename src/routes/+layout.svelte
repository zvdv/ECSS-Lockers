<script>
  import "../app.css";
</script>

<slot />

<style>
  :global(html) {
    height: 100%;
  }
  :global(body) {
    min-height: 100%;
  }
</style>
