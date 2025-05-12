import { renderComponent } from "$lib/components/ui/data-table";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./data-table-actions.svelte";
import Title from "./title.svelte";
import SeverityCell from "./severity-cell.svelte";

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
    accessorKey: "client",
    header: "Client",
  },
  {
    accessorKey: "severity",
    header: "Severity",
    cell: ({ getValue }) => renderComponent(SeverityCell, { value: getValue() }),
  },
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ row }) => {
      return renderComponent(Title, { title: row.original.title });
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
