<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import SideNavigation from "$lib/components/navigation/sideNavigation.svelte";
  import TopNavigation from "$lib/components/navigation/topNavigation.svelte";
  let { children } = $props();
</script>

<Sidebar.Provider style="--sidebar-width: 16rem;">
  <SideNavigation />
  <Sidebar.Inset>
    <TopNavigation />
    <main>
      {@render children?.()}
    </main>
  </Sidebar.Inset>
</Sidebar.Provider>
