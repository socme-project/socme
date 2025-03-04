<script lang="ts">
  import Badge from "$lib/components/ui/badge/badge.svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>Last alerts</Card.Title>
    <Card.Description>Critical/High alerts only</Card.Description>
  </Card.Header>
  <Card.Content class="flex flex-col gap-4">
    <div class="flex justify-between gap-2">
      <div class="flex flex-col gap-1">
        <p>
          <Badge class="bg-red-400 hover:bg-red-400">Critical</Badge>
          Suspicious email by unknown domain
        </p>
        <p class="text-muted-foreground">Client:ESIEA ID:10123</p>
      </div>
      <p class="text-muted-foreground text-sm">20/01/22 22:12:12</p>
    </div>
    <div class="flex justify-between gap-2">
      <div class="flex flex-col gap-1">
        <p>
          <Badge class="bg-red-400 hover:bg-red-400">Critical</Badge>
          Suspicious email by unknown domain
        </p>
        <p class="text-muted-foreground">Client:ESIEA ID:10123</p>
      </div>
      <p class="text-muted-foreground text-sm">20/01/22 22:12:12</p>
    </div>
  </Card.Content>
  <Card.Footer class="flex justify-between">
    <Button variant="ghost" class="w-full">See all alerts</Button>
  </Card.Footer>
</Card.Root>
