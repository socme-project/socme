<script lang="ts">
import * as Card from "$lib/components/ui/card";
import { slide } from "svelte/transition";
import { ChartNoAxesColumnIncreasing, UserCog } from "lucide-svelte";
import { onMount } from "svelte";
let { active = 0, disconnected = 0 } = $props();

let max = Math.max(active, disconnected);

let load = $state(false)
onMount(() => {
  load = true
})
</script>

<Card.Root>
  <Card.Header
    class="flex flex-row items-center justify-between space-y-0 pb-2"
  >
    <Card.Title class="text-sm font-medium">Artemis clients status</Card.Title>
    <ChartNoAxesColumnIncreasing class="text-muted-foreground h-4 w-4" />
  </Card.Header>

  <div class="w-full grid grid-cols-2 px-6 gap-5 h-[28px]">
    <p class="text-center text-muted-foreground">
      Disconnected: {disconnected}
    </p>
    <p class="text-center text-muted-foreground">Active: {active}</p>
  </div>

      <Card.Content
    class="flex w-full gap-2 h-[80px]  items-end pb-0 mt-4 mb-0"
  >
    {#if load}
      <div transition:slide={{delay: 500, duration: 500}}
        style={"height:" + (disconnected/max)*80 + "px;"}
        class="bg-red-400 flex-grow rounded-t-lg"></div>
      <div transition:slide={{delay: 600, duration: 500}} 
        style={"height:" + (active/max)*80 + "px;"}
        class="h-20 bg-green-400 flex-grow rounded-t-lg"></div>
    {/if}
  </Card.Content>
</Card.Root>
