<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { isActiveLink } from 'sv-router';
  import {
    LayoutDashboard,
    Folder,
    Tv,
    ShieldCheck,
    ShieldAlert,
    FolderCog,
    UserCog,
  } from "@lucide/svelte";
  import User from "./user.svelte";
  import { user } from "$lib/stores/user";

  // Menu items.
  var items = $state([
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
        // {
        //   title: "TV View",
        //   url: "/tv-view",
        //   icon: Tv,
        // },
        {
          title: "Cert-fr",
          url: "/dashboard/certfr",
          icon: ShieldCheck,
        },
      ],
    },
  ]);

  $effect(() => {
    if (items.find((item) => item.label === "Administration")) return;
    if ($user?.Role === "admin") {
      items.push({
        label: "Administration",
        items: [
          {
            title: "Users",
            url: "/dashboard/admin/users",
            icon: UserCog,
          },
          {
            title: "Clients",
            url: "/dashboard/admin/clients",
            icon: FolderCog,
          },
        ],
      });
    }
  });
</script>

<Sidebar.Root variant="inset">
  <Sidebar.Header>
    <div class="flex items-center justify-center pt-4 pb-2 space-x-2">
      <img
        src="/logo-white.png"
        alt="logo"
        width="40"
        class="hidden dark:block"
      />
      <img
        src="/logo-black.png"
        alt="logo"
        width="40"
        class="block dark:hidden"
      />
      <p class="text-xl font-black">SOCme</p>
    </div>
  </Sidebar.Header>

  <Sidebar.Content>
    {#each items as group (group.label)}
      <Sidebar.Group>
        <Sidebar.GroupLabel>{group.label}</Sidebar.GroupLabel>
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            {#each group.items as item (item.title)}
              <Sidebar.MenuItem>
                <Sidebar.MenuButton
                  class="transition-all"
                >
                  {#snippet child({ props })}
                    <a href={item.url} {...props}
                      use:isActiveLink={{ className: 'bg-accent' }}
                    >
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            {/each}
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    {/each}
  </Sidebar.Content>

  <Sidebar.Footer>
    <Sidebar.Menu>
      <Sidebar.MenuItem>
        <User />
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Footer>
</Sidebar.Root>
