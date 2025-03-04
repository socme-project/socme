<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import ChevronsUpDown from "lucide-svelte/icons/chevrons-up-down";
  import LogOut from "lucide-svelte/icons/log-out";
  import UserCog from "lucide-svelte/icons/user-cog";

  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import { useSidebar } from "$lib/components/ui/sidebar/index.js";
  import { RefreshCw } from "lucide-svelte";
  const sidebar = useSidebar();
  import { user } from "$lib/stores/user";
</script>

{#if $user}
  <DropdownMenu.Root>
    <DropdownMenu.Trigger>
      {#snippet child({ props })}
        <Sidebar.MenuButton
          size="lg"
          class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          {...props}
        >
          <Avatar.Root class="h-8 w-8 rounded-lg">
            <Avatar.Image src={$user.Avatar} alt={$user.Username + " avatar"} />
            <Avatar.Fallback class="rounded-lg"
              >{$user.Username}</Avatar.Fallback
            >
          </Avatar.Root>
          <div class="grid flex-1 text-left text-sm leading-tight">
            <span class="truncate font-semibold">{$user.Username}</span>
            <span class="truncate text-xs">{$user.GitHubID} ({$user.Role})</span
            >
          </div>
          <ChevronsUpDown class="ml-auto size-4" />
        </Sidebar.MenuButton>
      {/snippet}
    </DropdownMenu.Trigger>
    <DropdownMenu.Content
      class="w-[--bits-dropdown-menu-anchor-width] min-w-56 rounded-lg"
      side={sidebar.isMobile ? "bottom" : "right"}
      align="end"
      sideOffset={4}
    >
      <DropdownMenu.Label class="p-0 font-normal">
        <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
          <Avatar.Root class="h-8 w-8 rounded-lg">
            <Avatar.Image src={$user.Avatar} alt={$user.Username + " avatar"} />
          </Avatar.Root>
          <div class="grid flex-1 text-left text-sm leading-tight">
            <span class="truncate font-semibold">{$user.Username}</span>
            <span class="truncate text-xs">{$user.GitHubID}</span>
          </div>
        </div>
      </DropdownMenu.Label>
      <DropdownMenu.Separator />
      <DropdownMenu.Group>
        <a href="/dashboard/settings">
          <DropdownMenu.Item>
            <UserCog />
            Settings
          </DropdownMenu.Item>
        </a>
      </DropdownMenu.Group>
      <DropdownMenu.Separator />
      <a href="/auth/refresh">
        <DropdownMenu.Item>
          <RefreshCw />
          Refresh Session
        </DropdownMenu.Item>
      </a>
      <a href="/auth/logout">
        <DropdownMenu.Item>
          <LogOut />
          Log out
        </DropdownMenu.Item>
      </a>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
{/if}
