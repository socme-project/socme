import { renderComponent } from "$lib/components/ui/data-table";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./data-table-actions.svelte";
import Title from "./title.svelte";
import SeverityCell from "./severity-cell.svelte";
import type { Client } from "$src/lib/stores/client";

export type Alert = {
  ID: number;
  ClientID: string;
  WazuhAlertID: string;
  RuleDescription: string;
  RuleID: string;
  RuleLevel: number;
  Sort: number;
  Timestamp: string;
  RawJSON: string;
  Client: Client;
};

export const columns: ColumnDef<Alert>[] = [
  {
    accessorKey: "client",
    header: "Client",
  },
  {
    accessorKey: "severity",
    header: "Severity",
    cell: ({ getValue }) =>
      renderComponent(SeverityCell, { value: getValue() }),
  },
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ row }) => {
      return renderComponent(Title, { title: row.original.RuleDescription });
    },
  },
  {
    accessorKey: "timestamp",
    header: "Timestamp",
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return renderComponent(DataTableActions, { id: row.original.ID });
    },
  },
];
