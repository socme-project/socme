<script lang="ts">
  import { ShieldAlert } from "lucide-svelte";
  import { columns } from "./columns";
  import DataTable from "./data-table.svelte";
  import type { Alert } from "./columns";
  import * as Pagination from "$lib/components/ui/pagination/index.js";

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
</script>

<h1 class="flex items-center gap-4 mb-8">
  <ShieldAlert />
  Alerts
</h1>

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
