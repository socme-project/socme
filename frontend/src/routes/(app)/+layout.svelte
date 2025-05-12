<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import SideNavigation from "$lib/components/navigation/sideNavigation.svelte";
  import TopNavigation from "$lib/components/navigation/topNavigation.svelte";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { toast } from "svelte-sonner";
  import { user } from "$lib/stores/user";
  import axios from "axios";
  let { children } = $props();

  onMount(async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      toast.error("Unauthorized");
      user.set(null);
      await goto("/");
    }
    await axios
      .get("/api/auth/refresh", {
        headers: { Authorization: token },
      })
      .then((res) => {
        user.set(res.data.user);
      })
      .catch(async () => {
        toast.error("Unauthorized, try to login again");
        localStorage.removeItem("token");
        user.set(null);
        await goto("/");
      });
  });
</script>

<Sidebar.Provider style="--sidebar-width: 16rem;">
  <SideNavigation />
  <Sidebar.Inset>
    <TopNavigation />
    <main class="bg-background">
      {@render children?.()}
    </main>
  </Sidebar.Inset>
</Sidebar.Provider>
