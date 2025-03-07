<script lang="ts">
  import { page } from "$app/state";
  import { ShieldAlert } from "lucide-svelte";

  import type { Alert } from "../columns";

  let id = page.url.pathname.split("/").pop();

  let alert: Alert = {
    id: 1,
    title: "Alert 1",
    severity: "high",
    client: "esiea",
    timestamp: "2021-10-01T12:00:00Z",
    raw: JSON.stringify({
      id: 1,
      title: "Alert 1",
      severity: "high",
      client: "esiea",
      timestamp: "2021-10-01T12:00:00Z",
    }),
  };
</script>

<h1 class="flex items-center gap-4 mb-8">
  <ShieldAlert />
  Alert #{id}
</h1>

<h2>Alert details</h2>

<table class="w-full text-left">
  <tbody>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >id</th
      >
      <td>{alert.id}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Severity</th
      >
      <td>{alert.severity}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Description</th
      >
      <td>{alert.title}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Client</th
      >
      <td>{alert.client}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Time</th
      >
      <td>{alert.timestamp}</td>
    </tr>
  </tbody>
</table>

<h2>Raw data</h2>

<table class="w-full text-left">
  <tbody>
    {#each Object.entries(JSON.parse(alert.raw)) as [key, value]}
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >{key}</th
        >
        <td>{value}</td>
      </tr>
    {/each}
  </tbody>
</table>
