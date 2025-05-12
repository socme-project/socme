<script lang="ts">
  import Badge from "$lib/components/ui/badge/badge.svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import type { Alert } from "../alerts/columns";

  const { data }: { data: Alert[] } = $props();
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>Last alerts</Card.Title>
    <Card.Description>Critical/High alerts only</Card.Description>
  </Card.Header>
  <Card.Content class="flex flex-col gap-4">
    {#if data.length === 0}
      <p class="text-muted-foreground">No alerts found</p>
    {:else}
      {#each data as alert}
        <a
          href={"/dashboard/alerts/" + alert.id}
          class="hover:bg-accent/40 p-4 rounded-lg"
        >
          <div class="flex justify-between gap-2">
            <div class="flex flex-col gap-1">
              <p>
                <Badge class="bg-red-400 hover:bg-red-400">Critical</Badge>
                {alert.title}
              </p>
              <p class="text-muted-foreground">
                Client:{alert.client} ID:{alert.id}
              </p>
            </div>
            <p class="text-muted-foreground text-sm">{alert.timestamp}</p>
          </div>
        </a>
      {/each}
    {/if}
  </Card.Content>
  <Card.Footer class="flex justify-between">
    <Button href="/dashboard/alerts" variant="ghost" class="w-full"
      >See all alerts</Button
    >
  </Card.Footer>
</Card.Root>
