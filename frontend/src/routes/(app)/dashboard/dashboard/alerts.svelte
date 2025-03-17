<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Tooltip from "$lib/components/ui/tooltip/index.js";
  import { slide } from "svelte/transition";
  import { onMount } from "svelte";

  let { alerts = [], color = "#F87171", text = "Critical" } = $props();

  let max = Math.max(...alerts, 1);
  let alertsGraph: any = $state([]);
  let total = $state(0);

  // AlertGraph is an array of objects with a value property and a height (from 0 to 100) property.
  $effect(() => {
    total = alerts.reduce((acc, curr) => acc + curr, 0);
    alertsGraph = alerts.map((alert) => {
      if (alert === 0) {
        return {
          value: alert,
          height: 1,
        };
      }
      return {
        value: alert,
        height: (alert / max) * 100,
      };
    });
  });
</script>

<Card.Root>
  <Card.Header
    class="flex flex-row items-center justify-between space-y-0 pb-2"
  >
    <Card.Title class="text-sm font-medium">{text} alerts</Card.Title>
    <p class="font-black" style={"color:" + color + ";"}>
      {total}
    </p>
  </Card.Header>
  <Card.Content class="flex w-full gap-2 h-[100px] items-end pb-0 mt-4 mb-0">
    {#each alertsGraph as item, i}
      <div
        transition:slide={{ delay: 500 + i * 100, duration: 1500 }}
        class="h-[100px] flex-grow flex items-end hover:bg-accent rounded-t-sm"
      >
        <Tooltip.Provider>
          <Tooltip.Root>
            <Tooltip.Trigger
              class="h-[100px] flex-grow flex items-end hover:bg-accent rounded-t-sm"
            >
              <div
                style={"height:" +
                  item.height +
                  "px; background-color:" +
                  color +
                  ";"}
                class="flex-grow rounded-t-sm"
              ></div>
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p>{item.value} alerts</p>
            </Tooltip.Content>
          </Tooltip.Root>
        </Tooltip.Provider>
      </div>
    {/each}
  </Card.Content>
</Card.Root>
