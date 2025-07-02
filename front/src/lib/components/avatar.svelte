<script lang="ts">
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import type { User } from "../stores/user";
  import { cn } from "../utils";

  let {
    user,
    className,
  }: {
    user: User;
    className?: string;
  } = $props();

  const url =
    "https://avatars.githubusercontent.com/u/" + user.GithubID + "?v=4";
</script>

{#if user}
<Avatar.Root class={cn("h-8 w-8 rounded-lg", className)}>
  <Avatar.Image src={url} alt={user.Name + " avatar"} />
  <Avatar.Fallback class="rounded-lg">{user.Name}</Avatar.Fallback>
</Avatar.Root>
{:else}
<Avatar.Root class={cn("h-8 w-8 rounded-lg", className)}>
  <Avatar.Fallback class="rounded-lg">Unknown</Avatar.Fallback>
</Avatar.Root>
{/if}
