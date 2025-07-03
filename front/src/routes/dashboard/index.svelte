<script lang="ts">
  import * as Chart from "$lib/components/ui/chart/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { PieChart, Text } from "layerchart";
  import Alerts from "$src/lib/components/charts/alerts.svelte";
  import Button from "$src/lib/components/ui/button/button.svelte";
    import Clients from "$src/lib/components/charts/clients.svelte";

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
</script>

<h1>Dashboard</h1>

<div class="grid lg:grid-cols-4 sm:grid-cols-2 grid-cols-1 gap-5">
  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>Clients status</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1">
      <Clients actif={12} inactif={2}/>
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>Low severity</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={[3, 23, 20, 23, 12, 13]}
        hexColor="oklch(0.72 0.19 150)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>Medium severity</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={[3, 23, 20, 23, 12, 13]}
        hexColor="oklch(0.70 0.19 48)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center mb-6">
      <Card.Title>High severity</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={[3, 23, 20, 23, 12, 13]}
        hexColor="oklch(0.64 0.21 25)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col col-span-full">
    <Card.Header class="items-center">
      <Card.Title>Last alerts</Card.Title>
      <Card.Description>
        Only medium/high alerts displayed
      </Card.Description>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      ...
      ...
    </Card.Content>
    <Card.Footer>
      <Button href="/dashboard/alerts">See all alerts</Button>
    </Card.Footer>
  </Card.Root>
</div>
