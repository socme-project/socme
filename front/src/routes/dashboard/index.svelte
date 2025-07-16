<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import Alerts from "$src/lib/components/charts/alerts.svelte";
  import Button from "$src/lib/components/ui/button/button.svelte";
  import Clients from "$src/lib/components/charts/clients.svelte";
  import {
    SignalMedium,
    SignalHigh,
    Signal,
    LayoutDashboard,
    SignalLow,
  } from "@lucide/svelte";
  import { sendError } from "$src/lib/utils";
  import axios from "axios";
  import { onMount } from "svelte";
  import Badge from "$src/lib/components/ui/badge/badge.svelte";
  import type { Alert } from "$src/lib/components/alerts/columns";

  let statsLow = $state([]);
  let statsMedium = $state([]);
  let statsHigh = $state([]);
  let statsCritical = $state([]);
  let lastAlerts = $state<Alert[]>([]);
  let activeAgents = $state<number>(0);
  let inactiveAgents = $state<number>(0);
  let activeClients = $state<number>(0);
  let inactiveClients = $state<number>(0);

  async function loadStats() {
    axios
      .get("/api/alerts/stats/agents")
      .then((response) => {
        activeAgents = response.data.activeAgents;
        inactiveAgents = response.data.inactiveAgents;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios
      .get("/api/alerts/stats/clients")
      .then((response) => {
        activeClients = response.data.activeClients;
        inactiveClients = response.data.inactiveClients;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios
      .get("/api/alerts/stats/low")
      .then((response) => {
        statsLow = response.data.events;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios
      .get("/api/alerts/stats/medium")
      .then((response) => {
        statsMedium = response.data.events;
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios
      .get("/api/alerts/stats/high")
      .then((response) => {
        statsHigh = response.data.events;
        console.log(response.data.events);
      })
      .catch((error) => {
        sendError("Error fetching stats:", error);
      });
    axios
      .get("/api/alerts/stats/critical")
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
        lastAlerts = res.data.alerts;
      })
      .catch((err) => {
        sendError("Error fetching stats:", err);
      });
  }

  onMount(async () => {
    loadStats();
    setInterval(
      () => {
        loadStats();
      },
      5 * 60 * 1000,
    );
  });
</script>

<h1 class="flex items-center gap-4">
  <LayoutDashboard />
  Dashboard
</h1>

<div class="grid xl:grid-cols-4 lg:grid-cols-2 grid-cols-1 gap-5">
  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>Clients status</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1">
      <Clients actif={activeClients} inactif={inactiveClients} />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>Total agents status</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1">
      <Clients actif={activeAgents} inactif={inactiveAgents} />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><SignalLow /> Low severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      {#if statsLow.length === 0}
        <Alerts alerts={statsLow} hexColor="oklch(0.70 0.19 48)" />
      {:else}
        <Alerts alerts={statsLow} hexColor="oklch(0.70 0.19 48)" />
      {/if}
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><SignalMedium /> Medium severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      {#if statsMedium.length === 0}
        <Alerts alerts={statsMedium} hexColor="oklch(0.70 0.19 48)" />
      {:else}
        <Alerts alerts={statsMedium} hexColor="oklch(0.70 0.19 48)" />
      {/if}
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><SignalHigh /> High severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      {#if statsHigh.length === 0}
        <Alerts alerts={statsHigh} hexColor="oklch(0.70 0.19 48)" />
      {:else}
        <Alerts alerts={statsHigh} hexColor="oklch(0.70 0.19 48)" />
      {/if}
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title class="flex gap-1 items-center "
        ><Signal /> Critical severity</Card.Title
      >
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      {#if statsCritical.length === 0}
        <Alerts alerts={statsCritical} hexColor="oklch(0.70 0.19 48)" />
      {:else}
        <Alerts alerts={statsCritical} hexColor="oklch(0.70 0.19 48)" />
      {/if}
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
