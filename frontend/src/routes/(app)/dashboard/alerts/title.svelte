<script lang="ts">
  import Badge from "$lib/components/ui/badge/badge.svelte";

  let { title, severity } = $props();
</script>

<div class="flex items-center gap-2">
  {#if severity === "critical"}
    <Badge variant="destructive">{severity}</Badge>
  {:else if severity === "high"}
    <Badge
      class="bg-orange-400 hover:bg-orange-400/80 dark:bg-orange-700 hover:bg-orange-700/80"
      >{severity}</Badge
    >
  {:else}
    <Badge variant="outline">{severity}</Badge>
  {/if}
  {title}
</div>
