<script lang="ts">
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";
  import * as Card from "$lib/components/ui/card/index.js";
</script>

<Card.Root>
  <Card.Header>
    <Card.Title class="font-medium">
      <Skeleton class="w-3/5 h-8 mb-1" />
      <Skeleton class="w-2/3 h-4" />
    </Card.Title>
  </Card.Header>
  <Card.Content class="flex flex-col gap-1 opacity-70">
    <Skeleton class="w-4/5 h-4" />
    <Skeleton class="w-5/5 h-4" />
    <Skeleton class="w-3/5 h-4" />
  </Card.Content>
</Card.Root>
