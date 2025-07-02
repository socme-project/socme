<script lang="ts">
  import { useSidebar } from "$lib/components/ui/sidebar/index.js";
  import { user } from "$lib/stores/user";
  import { ChevronsUpDown, LogOut, RefreshCw, UserCog } from "@lucide/svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import Avatar from "$lib/components/avatar.svelte"
  const sidebar = useSidebar();
</script>

{#if $user}
  <DropdownMenu.Root>
    <DropdownMenu.Trigger>
      {#snippet child({ props }: { props: any })}
        <Sidebar.MenuButton
          size="lg"
          class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          {...props}
        >
          <Avatar user={$user} />
          <div class="grid flex-1 text-left text-sm leading-tight">
            <span class="truncate font-semibold">{$user.Name}</span>
            <span class="truncate text-xs">{$user.GithubID} ({$user.Role})</span
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
          <Avatar user={$user} />
          <div class="grid flex-1 text-left text-sm leading-tight">
            <span class="truncate font-semibold">{$user.Name}</span>
            <span class="truncate text-xs">{$user.GithubID}</span>
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
{:else}
  <p class="w-full text-center text-sm text-muted-foreground">Disconnected</p>
{/if}
