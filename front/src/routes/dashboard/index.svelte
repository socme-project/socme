<script lang="ts">
  import * as Chart from "$lib/components/ui/chart/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { PieChart, Text } from "layerchart";
  import Alerts from "$src/lib/components/alerts.svelte";
  import Button from "$src/lib/components/ui/button/button.svelte";

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

<div class="grid grid-cols-4 gap-5">
  <Card.Root class="flex flex-col">
    <Card.Header class="items-center">
      <Card.Title>Statut des agents Wazuh</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1">
      <Chart.Container
        config={chartConfig}
        class="aspect-square max-w-40 mx-auto"
      >
        <PieChart
          data={agentData}
          key="status"
          value="count"
          c="color"
          innerRadius={60}
          padding={28}
          props={{ pie: { motion: "tween" } }}
        >
          {#snippet aboveMarks()}
            <Text
              value={String(totalAgents)}
              textAnchor="middle"
              verticalAnchor="middle"
              class="fill-foreground text-3xl! font-bold"
              dy={3}
            />
            <Text
              value="Agents"
              textAnchor="middle"
              verticalAnchor="middle"
              class="fill-muted-foreground! text-muted-foreground"
              dy={22}
            />
          {/snippet}
          {#snippet tooltip()}
            <Chart.Tooltip hideLabel />
          {/snippet}
        </PieChart>
      </Chart.Container>
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center">
      <Card.Title>Low severity</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={[10, 2, 3, 4, 5, 6, 23, 23, 20, 23, 12, 13]}
        hexColor="oklch(0.72 0.19 150)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center">
      <Card.Title>Medium severity</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={[10, 2, 3, 4, 5, 6, 23, 23, 20, 23, 12, 13]}
        hexColor="oklch(0.70 0.19 48)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col">
    <Card.Header class="items-center">
      <Card.Title>High severity</Card.Title>
    </Card.Header>
    <Card.Content class="flex-1 grid items-end">
      <Alerts
        alerts={[10, 2, 3, 4, 5, 6, 23, 23, 20, 23, 12, 13]}
        hexColor="oklch(0.64 0.21 25)"
      />
    </Card.Content>
  </Card.Root>

  <Card.Root class="flex flex-col col-span-4">
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
      <Button>See all alerts</Button>
    </Card.Footer>
  </Card.Root>
</div>
