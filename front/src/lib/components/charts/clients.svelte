<script lang="ts">
  import { PieChart } from "layerchart";
  import * as Chart from "$lib/components/ui/chart/index.js";

  let { actif, inactif } = $props();

  const chartData = [
    { status: "actif", count: actif, color: "var(--color-green-500)" },
    { status: "inactif", count: inactif, color: "var(--color-red-500)" },
  ];

  const chartConfig = {
    count: { label: "Count" },
    actif: { label: "Actif", color: "var(--chart-1)" },
    inactif: { label: "Inactif", color: "var(--chart-2)" },
  } satisfies Chart.ChartConfig;
</script>

{#if actif === 0 && inactif === 0}
  <p class="text-muted-foreground text-center">No data available</p>
{:else}
  <Chart.Container config={chartConfig} class="mx-auto aspect-square w-1/2">
    <PieChart
      data={chartData}
      key="status"
      value="count"
      cRange={chartData.map((d) => d.color)}
      c="color"
      props={{
        pie: {
          motion: "tween",
        },
      }}
    >
      {#snippet tooltip()}
        <Chart.Tooltip hideLabel />
      {/snippet}
    </PieChart>
  </Chart.Container>
{/if}
