<script lang="ts">
  import { invalidate } from "$app/navigation";
  import Badge from "./Badge.svelte";
  import Button from "./Button.svelte";
  import FilterIcon from "./FilterIcon.svelte";
  import FilterMenu from "./FilterMenu.svelte";

  export let data;

  let checked: {
    [k: string]: boolean;
  } = Object.fromEntries(
    data.registration.map(({ locker }) => [locker, false])
  );

  $: checkedLockers = Object.entries(checked)
    .filter(([, checked]) => checked)
    .map(([locker]) => locker);

  let showStatusFilter = false;
  let statusFilter = {
    available: true,
    claimed: true,
    expired: true,
  };

  $: filteredData = {
    registration: data.registration.filter((row) => statusFilter[row.status]),
  };

  function selectAll() {
    for (const locker of Object.keys(checked)) {
      checked[locker] = true;
    }
  }
  function deselectAll() {
    for (const locker of Object.keys(checked)) {
      checked[locker] = false;
    }
  }
  async function deleteLockers(lockers: string[]) {
    await fetch("/admin/api/lockers", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(lockers),
    });
    invalidate("/admin/api/lockers");
    deselectAll();
  }
  async function renew(locker: string) {
    await fetch("admin/renew", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ locker }),
    });
    invalidate("/admin/api/lockers");
  }
</script>

<div>
  <h1 class="p-2 text-3xl font-semibold">Locker Admin Panel</h1>
  <hr class="border-black" />
  <div class="p-2">
    {#if Object.values(checked).some((checked) => !checked)}
      <Button on:click={selectAll}>Select all</Button>
    {:else}
      <Button on:click={deselectAll}>Deselect all</Button>
    {/if}
    {#if Object.values(checked).some((checked) => checked)}
      <Button on:click={() => deleteLockers(checkedLockers)}
        >Delete selected</Button
      >
    {/if}
    <a href="admin/import"><Button>Import from Excel</Button></a>
    <a href="admin/export" download><Button>Export to Excel</Button></a>
  </div>
  <table class="w-full border-collapse border-spacing-0">
    <tr class="bg-neutral-200">
      <th />
      <th>Locker</th>
      <th>
        <Button on:click={() => (showStatusFilter = !showStatusFilter)}>
          Status
          <FilterIcon class="inline-block h-4 w-4" />
        </Button>
        <FilterMenu hidden={!showStatusFilter} bind:filter={statusFilter} />
      </th>
      <th>Name</th>
      <th>Email</th>
      <th>Expiry</th>
      <th />
    </tr>
    {#each filteredData.registration as { locker, status, name, user, expiry } (locker)}
      <tr class={status} class:checked={checked[locker]}>
        <td class="w-0 text-center">
          <div class="flex h-full w-full items-center justify-center">
            <input
              type="checkbox"
              bind:checked={checked[locker]}
              class="h-5 w-5"
            />
          </div>
        </td>
        <td>{locker}</td>
        <td><Badge {status}>{status}</Badge></td>
        <td>{name ?? ""}</td>
        <td>{user ?? ""}</td>
        <td>{expiry?.toDateString() ?? ""}</td>
        <td>
          {#if status === "available"}
            <a href="admin/register?locker={locker}"
              ><Button>Register</Button></a
            >
          {:else}
            <a href="admin/edit?locker={locker}"><Button>Edit</Button></a>
            <a href="admin/delete?locker={locker}"><Button>Delete</Button></a>
          {/if}
          {#if status === "expired"}
            <Button on:click={() => renew(locker)}>Renew</Button>
          {/if}
        </td>
      </tr>
    {/each}
  </table>
</div>

<style lang="postcss">
  th {
    border-top: 1px black solid;
    position: sticky;
    top: 0;
    @apply bg-neutral-200;
  }
  th,
  td {
    border: 1px black solid;
    border-right: none;
    @apply p-1;
  }
  td {
    border-top: none;
    @apply py-0.5;
  }

  .available {
    @apply bg-neutral-100;
  }
  .claimed {
    @apply bg-green-100;
  }
  .expired {
    @apply bg-red-100;
  }
  .checked {
    @apply bg-sky-100;
  }
</style>
