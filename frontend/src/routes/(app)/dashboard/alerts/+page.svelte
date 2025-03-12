<script lang="ts">
  import {
    FileQuestionIcon,
    ShieldAlert,
    TouchpadOffIcon,
  } from "lucide-svelte";
  import { columns } from "./columns";
  import DataTable from "./data-table.svelte";
  import type { Alert } from "./columns";
  import * as Pagination from "$lib/components/ui/pagination/index.js";
  import Filters from "./filters.svelte";
  import FiltersSingle from "./filtersSingle.svelte";
  import Input from "$lib/components/ui/input/input.svelte";

  let currentPage = $state(1);

  let perPage = $state(10);
  let count = $state(100);

  let alerts: Alert[] = [
    {
      id: 1,
      title: "Alert 1",
      severity: "high",
      client: "esiea",
      timestamp: "2021-10-01T12:00:00Z",
    },
  ];

  export const statuses = [
    {
      value: "backlog",
      label: "Backlog",
      icon: FileQuestionIcon,
    },
    {
      value: "todo",
      label: "Todo",
      icon: TouchpadOffIcon,
    },
  ];

  let selectedValues = $state([]);
  let selectedValue = $state("");
  let search = $state("");
</script>

<h1 class="flex items-center gap-4 mb-8">
  <ShieldAlert />
  Alerts
</h1>

<div class=" my-10 flex gap-5 flex-wrap">
  <Input
    bind:value={search}
    placeholder="Filter tasks..."
    class="h-8 w-[150px] lg:w-[250px]"
  />
  <p class="text-muted">Filters:</p>
  <Filters title="Severity" options={statuses} bind:selectedValues />
  <FiltersSingle title="Clients" options={statuses} bind:selectedValue />
  <Filters title="Tag" options={statuses} bind:selectedValues />
  <!-- Rule ID // maybe in search -->
  <!-- Time // maybe in search -->
</div>

<DataTable data={alerts} {columns} />

<Pagination.Root
  {count}
  {perPage}
  class="mt-6"
  onPageChange={(p) => {
    currentPage = p;
  }}
>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.PrevButton />
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link {page} isActive={currentPage === page.value}>
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
