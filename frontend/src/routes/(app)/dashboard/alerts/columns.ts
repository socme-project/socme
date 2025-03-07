import { renderComponent, renderSnippet } from "$lib/components/ui/data-table";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./data-table-actions.svelte";
import Title from "./title.svelte";

export type Alert = {
  id: number;
  title: string;
  severity: "critical" | "high" | "medium" | "low";
  timestamp: string;
  client: string;
  raw: string;
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
    cell: ({ row }) => {
      return renderComponent(Title, { title: row.original.title, severity: row.original.severity });
    },
  },
  {
    accessorKey: "timestamp",
    header: "Timestamp",
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return renderComponent(DataTableActions, { id: row.original.id });
    },
  }
];
