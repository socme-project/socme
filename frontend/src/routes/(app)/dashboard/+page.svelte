<script lang="ts">
  import ArtemisStatus from "./dashboard/artemis-status.svelte";
  import Alerts from "./dashboard/alerts.svelte";
  import { LayoutDashboard } from "lucide-svelte";
  import axios from "axios";
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import type { Alert } from "./alerts/columns";
  import AlertList from "./dashboard/alert-list.svelte";

  let lastFiveAlerts: Alert[] = $state([]);

  onMount(() => {
    axios
      .get("/api/alerts/getlastfive", {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        lastFiveAlerts = res.data.alerts.map((alert: any) => ({
          id: alert.ID,
          title: alert.rule_description,
          severity: alert.rule_level,
          client: alert.client_name,
          timestamp: alert.timestamp,
        }));
      })
      .catch(() => {
        toast.error("Internal server error");
      });

    axios
      .get("/api/alerts/last24h/high", {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        console.log(res.data);
      })
      .catch(() => {
        toast.error("Internal server error");
      });
  });
</script>

<h1 class="flex items-center gap-4 mb-8">
  <LayoutDashboard />
  Dashboard
</h1>
<div class="flex flex-col gap-4">
  <div class="grid gap-4 grid-cols-1 sm:grid-cols-2 xl:grid-cols-4">
    <ArtemisStatus active={21} disconnected={2} />
    <Alerts color="#F87171" text="Critical" alerts={[1, 2, 3, 2, 5, 4, 4, 3]} />
    <Alerts color="#FB923C" text="High" alerts={[100, 50, 3, 2, 20, 4, 4, 3]} />
    <Alerts
      color="#60A5FA"
      text="Medium"
      alerts={[50, 50, 63, 80, 60, 40, 40, 30]}
    />
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2">
    <AlertList data={lastFiveAlerts} />
  </div>
</div>
