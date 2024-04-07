<script lang="ts">
  import chevron from "./chevron_right.svg";

  export let rows: Row[];

  type Row = {
    name: string;
    locker: string;
  };
</script>

<div class="wrapper">
  <table>
    {#each rows as row}
      <tr class="relative">
        <!-- stretch a link across the entire table row -->
        <!-- this works using "relative" inside "absolute inset-0" -->
        <!-- i put the div there to silence accessibility warnings (yes, i know) -->
        <a class="absolute inset-0" href="lockers/{row.locker}"><div /></a>

        <td>{row.name}</td>
        <td>{row.locker}</td>
        <td>
          <img class="ml-auto h-7 w-7" src={chevron} alt="Right chevron" />
        </td>
      </tr>
    {/each}
  </table>
</div>

<style lang="postcss">
  .wrapper {
    @apply w-full max-w-sm bg-neutral-50;
    @apply rounded-lg border border-neutral-300;
    /* we round table by wrapping it in a rounded div and hiding overflow */
    @apply box-border overflow-hidden;
  }
  table {
    width: 100%;
    border-collapse: collapse;
  }
  tr {
    @apply text-lg hover:bg-neutral-200;
  }
  tr:last-child td {
    border-bottom: 0px;
  }
  td {
    @apply p-3;
    border-bottom: 1px solid theme(colors.neutral.300);
  }
</style>
