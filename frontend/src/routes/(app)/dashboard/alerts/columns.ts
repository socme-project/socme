import { renderComponent, renderSnippet } from "$lib/components/ui/data-table";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./data-table-actions.svelte";
import { createRawSnippet } from "svelte";

export type Alert = {
  id: number;
  title: string;
  severity: "critical" | "high" | "medium" | "low";
  timestamp: string;
  client: string;
};

export const columns: ColumnDef<Alert>[] = [
  {
    accessorKey: "id",
    header: "ID",
  },
  {
    accessorKey: "client",
    header: "Client",
  },
  {
    accessorKey: "title",
    header: "Title",
  },
  {
    accessorKey: "severity",
    header: "Severity",
  },
  {
    accessorKey: "timestamp",
    header: "Timestamp",
  },
  {
    id: "actions",
    cell: ({ row }) => {
      // You can pass whatever you need from `row.original` to the component
      return renderComponent(DataTableActions, { id: row.original.id });
    },
  }
];
