<script lang="ts">
  import * as Command from "$lib/components/ui/command/index.js";
  import { cn } from "$src/lib/utils";
  import Button from "../ui/button/button.svelte";
  import {
    CloudMoon,
    GithubIcon,
    LogOut,
    Moon,
    RefreshCw,
    Settings,
    ShieldCheck,
    Sun,
    LayoutDashboard,
    ShieldAlert,
    Folder,
    Tv,
  } from "@lucide/svelte";
  import { resetMode, setMode } from "mode-watcher";

  let open = $state(false);

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
      e.preventDefault();
      open = !open;
    }
  }

  var groups = [
    {
      label: "Application",
      items: [
        {
          title: "Dashboard",
          url: "/dashboard",
          icon: LayoutDashboard,
        },
        {
          title: "Alerts",
          url: "/dashboard/alerts",
          icon: ShieldAlert,
        },
        {
          title: "Clients",
          url: "/dashboard/clients",
          icon: Folder,
        },
        {
          title: "Cert-fr",
          url: "/dashboard/certfr",
          icon: ShieldCheck,
        },
      ],
    },
  ];
</script>

<svelte:document onkeydown={handleKeydown} />

<Button
  onclick={() => (open = true)}
  class={cn(
    "focus-visible:ring-ring inline-flex justify-between items-center gap-2 whitespace-nowrap rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground border shadow-sm h-9 px-4 py-2 text-muted-foreground relative text-sm",
    "flex-grow",
  )}
>
  Search
  <kbd
    class="bg-muted pointer-events-none hidden select-none items-center gap-1 rounded border px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex"
  >
    <span class="text-xs">ctrl</span>K
  </kbd>
</Button>

<Command.Dialog bind:open>
  <Command.Input placeholder="Type a command or search..." />
  <Command.List>
    <Command.Empty>No results found.</Command.Empty>

    {#each groups as group}
      <Command.Group heading={group.label}>
        {#each group.items as item}
          <Command.LinkItem href={item.url} onSelect={() => (open = !open)}>
            <item.icon class="mr-2 size-4" />
            <span>{item.title}</span>
          </Command.LinkItem>
        {/each}
      </Command.Group>
    {/each}

    <Command.Separator />
    <Command.Group heading="Other">
      <Command.LinkItem
        onSelect={() => (open = !open)}
        href="/dashboard/settings"
      >
        <Settings class="mr-2 size-4" />
        <span>User settings</span>
      </Command.LinkItem>
      <Command.LinkItem onSelect={() => (open = !open)} href="/auth/logout">
        <LogOut class="mr-2 size-4" />
        <span>Log out</span>
      </Command.LinkItem>
      <Command.LinkItem onSelect={() => (open = !open)} href="/auth/refresh">
        <RefreshCw class="mr-2 size-4" />
        <span>Refresh session</span>
      </Command.LinkItem>
      <Command.LinkItem
        onSelect={() => (open = !open)}
        href="https://github.com/socme-project"
        target="_blank"
      >
        <GithubIcon class="mr-2 size-4" />
        <span>Open the github repo</span>
      </Command.LinkItem>

      <Command.Item
        onSelect={() => {
          open = !open;
          setMode("dark");
        }}
      >
        <Moon class="mr-2 size-4" />
        <span>Dark theme</span>
      </Command.Item>

      <Command.Item
        onSelect={() => {
          open = !open;
          setMode("light");
        }}
      >
        <Sun class="mr-2 size-4" />
        <span>Light theme</span>
      </Command.Item>

      <Command.Item
        onSelect={() => {
          open = !open;
          resetMode();
        }}
      >
        <CloudMoon class="mr-2 size-4" />
        <span>System theme</span>
      </Command.Item>
    </Command.Group>
  </Command.List>
</Command.Dialog>
