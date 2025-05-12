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

  let highPerHour: number[] = $state([]);
  let mediumPerHour: number[] = $state([]);
  let criticalPerHour: number[] = $state([]);

  onMount(() => {
    axios
      .get("/api/alerts", {
        headers: { Authorization: localStorage.getItem("token") },
        params: {
          page: 1,
          perPage: 5,
        },
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
      .get("/api/alerts/stats/high", {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        highPerHour = res.data.events;
      })
      .catch(() => {
        toast.error("Internal server error");
      });

    axios
      .get("/api/alerts/stats/critical", {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        criticalPerHour = res.data.events;
      })
      .catch(() => {
        toast.error("Internal server error");
      });

    axios
      .get("/api/alerts/stats/medium", {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        mediumPerHour = res.data.events;
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
    <Alerts color="#F87171" text="Critical" alerts={criticalPerHour} />
    <Alerts color="#FB923C" text="High" alerts={highPerHour} />
    <Alerts color="#60A5FA" text="Medium" alerts={mediumPerHour} />
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2">
    <AlertList data={lastFiveAlerts} />
  </div>
</div>
