<script lang="ts" generics="TData, TValue">
  import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";
  import {
    createSvelteTable,
    FlexRender,
  } from "$lib/components/ui/data-table/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as HoverCard from "$lib/components/ui/hover-card/index.js";
  type DataTableProps<TData, TValue> = {
    columns: ColumnDef<TData, TValue>[];
    data: TData[];
  };
  import CalendarDays from "lucide-svelte/icons/calendar-days";
  import * as Avatar from "$lib/components/ui/avatar/index.js";

  let { data, columns }: DataTableProps<TData, TValue> = $props();

  const table = createSvelteTable({
    get data() {
      return data;
    },
    columns,
    getCoreRowModel: getCoreRowModel(),
  });
</script>

<div class="rounded-md border">
  <Table.Root>
    <Table.Header>
      {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
        <Table.Row>
          {#each headerGroup.headers as header (header.id)}
            <Table.Head>
              {#if !header.isPlaceholder}
                <FlexRender
                  content={header.column.columnDef.header}
                  context={header.getContext()}
                />
              {/if}
            </Table.Head>
          {/each}
        </Table.Row>
      {/each}
    </Table.Header>
    <Table.Body>
      {#each table.getRowModel().rows as row (row.id)}
        <Table.Row data-state={row.getIsSelected() && "selected"}>
          {#each row.getVisibleCells() as cell (cell.id)}
            <Table.Cell>
              {#if cell.column.id === "client"}
                <HoverCard.Root>
                  <HoverCard.Trigger
                    href="https://github.com/sveltejs"
                    target="_blank"
                    rel="noreferrer noopener"
                    class="rounded-sm underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-offset-8 focus-visible:outline-black"
                  >
                    {cell.row.original.client}
                  </HoverCard.Trigger>
                  <HoverCard.Content class="w-80">
                    <div class="flex justify-between space-x-4">
                      <Avatar.Root>
                        <Avatar.Image src="https://github.com/sveltejs.png" />
                        <Avatar.Fallback>SK</Avatar.Fallback>
                      </Avatar.Root>
                      <div class="space-y-1">
                        <h4 class="text-sm font-semibold">@sveltejs</h4>
                        <p class="text-sm">Cybernetically enhanced web apps.</p>
                        <div class="flex items-center pt-2">
                          <CalendarDays class="mr-2 size-4 opacity-70" />
                          <span class="text-muted-foreground text-xs">
                            Joined September 2022
                          </span>
                        </div>
                      </div>
                    </div>
                  </HoverCard.Content>
                </HoverCard.Root>
              {:else}
                <FlexRender
                  content={cell.column.columnDef.cell}
                  context={cell.getContext()}
                />
              {/if}
            </Table.Cell>
          {/each}
        </Table.Row>
      {:else}
        <Table.Row>
          <Table.Cell colspan={columns.length} class="h-24 text-center">
            No results.
          </Table.Cell>
        </Table.Row>
      {/each}
    </Table.Body>
  </Table.Root>
</div>
