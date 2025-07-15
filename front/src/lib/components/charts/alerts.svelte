<script lang="ts">
  import { scaleBand } from "d3-scale";
  import { BarChart, type ChartContextValue } from "layerchart";
  import * as Chart from "$lib/components/ui/chart/index.js";
  import { cubicInOut } from "svelte/easing";

  const {
    alerts,
    hexColor,
  }: { alerts: number[]; hexColor: string } = $props();

  const chartData = [
    { hour: "-20h", count: alerts[0] },
    { hour: "-16h", count: alerts[1] },
    { hour: "-12h", count: alerts[2] },
    { hour: "-8h", count: alerts[3] },
    { hour: "-4h", count: alerts[4] },
    { hour: "now", count: alerts[5] },
  ];

  const chartConfig = {
    hour: { label: "Hour", color: hexColor },
  } satisfies Chart.ChartConfig;

  let context = $state<ChartContextValue>();

  // TODO: Float to int
</script>

<Chart.Container config={chartConfig}>
  <BarChart
    labels={{ offset: 12 }}
    data={chartData}
    xScale={scaleBand().padding(0.25)}
    x="hour"
    series={[
      { key: "count", label: "Alerts", color: chartConfig.hour.color },
    ]}
    axis="x"
    rule={false}
    props={{
      bars: {
        stroke: "none",
        radius: 8,
        rounded: "all",
        // use the height of the chart to animate the bars
        initialY: (context?.height ?? 0) + 180,
        initialHeight: 0,
        motion: {
          y: { type: "tween", duration: 500, easing: cubicInOut },
          height: { type: "tween", duration: 500, easing: cubicInOut },
        },
      },
      highlight: { area: { fill: "none" } },
      xAxis: { format: (d) => d.slice(0, 3) },
    }}
  >
    {#snippet tooltip()}
      <Chart.Tooltip hideLabel />
    {/snippet}
  </BarChart>
</Chart.Container>
