<script lang="ts">
  import * as Chart from "$lib/components/ui/chart/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { PieChart, Text } from "layerchart";
  import Alerts from "$src/lib/components/charts/alerts.svelte";
  import Button from "$src/lib/components/ui/button/button.svelte";
  import Clients from "$src/lib/components/charts/clients.svelte";
  import {
    SignalMedium,
    SignalHigh,
    Signal,
    LayoutDashboard,
  } from "@lucide/svelte";
    import { sendError } from "$src/lib/utils";
    import axios from "axios";
    import { onMount } from "svelte";
    import Badge from "$src/lib/components/ui/badge/badge.svelte";
    import type { Alert } from "$src/lib/components/alerts/columns";

  // TODO: Refresh every 5 minutes
  const agentData = [
    { status: "actif", count: 45, color: "var(--color-green-500)" },
    { status: "inactif", count: 12, color: "var(--color-red-500)" },
  ];

  const chartConfig = {
    count: { label: "Nombre d'agents" },
    actif: { label: "Actif", color: "var(--chart-1)" },
    inactif: { label: "Inactif", color: "var(--chart-2)" },
  } satisfies Chart.ChartConfig;

  const totalAgents = agentData.reduce((acc, curr) => acc + curr.count, 0);

  let statsMedium = $state([]);
  let statsHigh = $state([]);
  let statsCritical = $state([]);
  let lastAlerts = $state<Alert[]>([]);

  async function loadStats() {
    axios.get("/api/alerts/stats/medium")
      .then((response) => {
        statsMedium = response.data.events;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios.get("/api/alerts/stats/high")
      .then((response) => {
        statsHigh = response.data.events;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios.get("/api/alerts/stats/critical")
      .then((response) => {
        statsCritical = response.data.events;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });

    axios
      .get("/api/alerts", {
        params: {
          page: "1",
          perPage: 5,
          severity: "high,critical",
          preload: "true",
        },
      })
      .then((res) => {
        lastAlerts = res.data.alerts
      })
      .catch((err) => {
        sendError("Error fetching stats:", err);
      });
  }

  onMount(async ()=>{loadStats()})
</script>

<h1 class="flex items-center gap-4">
  <LayoutDashboard />
  Dashboard
</h1>

<div class="grid lg:grid-cols-4 sm:grid-cols-2 grid-cols-1 gap-5">
  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>Clients status</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1">
      <Clients actif={12} inactif={2} />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><SignalMedium /> Medium severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={statsMedium}
        hexColor="oklch(0.72 0.19 150)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><SignalHigh /> High severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts alerts={statsHigh} hexColor="oklch(0.70 0.19 48)" />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><Signal /> Critical severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts alerts={statsCritical} hexColor="oklch(0.64 0.21 25)" />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col col-span-full">
    <Card.Header class="items-center">
      <Card.Title>Last alerts</Card.Title>
      <Card.Description>Only high/critical alerts displayed</Card.Description>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
    {#if lastAlerts.length === 0}
      <p class="text-muted-foreground">No alerts found</p>
    {:else}
      {#each lastAlerts as alert}
        <a
          href={"/dashboard/alerts/" + alert.ID}
          class="hover:bg-accent/40 p-4 rounded-lg"
        >
          <div class="flex justify-between gap-2">
            <div class="flex flex-col gap-1">
              <p>
                <Badge class="bg-red-400 hover:bg-red-400">Critical</Badge>
                {alert.RuleDescription}
              </p>
              <p class="text-muted-foreground">
                Client:{alert.Client.Name} ID:{alert.ID}
              </p>
            </div>
            <p class="text-muted-foreground text-sm">{alert.Timestamp}</p>
          </div>
        </a>
      {/each}
    {/if}
    </Card.Content>
    <Card.Footer>
      <Button href="/dashboard/alerts">See all alerts</Button>
    </Card.Footer>
  </Card.Root>
</div>
