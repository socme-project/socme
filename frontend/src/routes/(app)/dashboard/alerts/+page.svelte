<script lang="ts">
  import {
    ChevronLastIcon,
    FileQuestionIcon,
    ShieldAlert,
    TouchpadOffIcon,
  } from "lucide-svelte";
  import { columns } from "./columns";
  import DataTable from "./data-table.svelte";
  import type { Alert } from "./columns";
  import * as Pagination from "$lib/components/ui/pagination/index.js";
  import Filters from "./filters.svelte";
  import FiltersSingle from "./filtersSingle.svelte";
  import Input from "$lib/components/ui/input/input.svelte";
  import { onMount } from "svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import Button from "$lib/components/ui/button/button.svelte";
  import { Next, Previous } from "$lib/components/ui/carousel";

  let currentPage = $state(1);

  let perPage = $state(10);
  let maxPage = $state(1);

  let alerts: Alert[] = $state([]);

  $effect(async () => {
    await axios
      .get("/api/alerts/page", {
        headers: { Authorization: localStorage.getItem("token") },
        params: { page: currentPage, perPage: perPage },
      })
      .then((res) => {
        alerts = res.data.alerts.map((alert: any) => ({
          id: alert.ID,
          title: alert.rule_description,
          severity: alert.rule_level,
          client: alert.client_name,
          timestamp: alert.timestamp,
          raw: alert.raw_json,
        }));
        maxPage = res.data.maxPage;
        console.log(res.data);
      })
      .catch(() => {
        toast.error("Internal server error");
      });
  });

  export const statuses = [
    {
      value: "backlog",
      label: "Backlog",
      icon: FileQuestionIcon,
    },
    {
      value: "todo",
      label: "Todo",
      icon: TouchpadOffIcon,
    },
  ];

  let selectedValues = $state([]);
  let selectedValue = $state("");
  let search = $state("");
</script>

<h1 class="flex items-center gap-4 mb-8">
  <ShieldAlert />
  Alerts
</h1>

<div class=" my-10 flex gap-5 flex-wrap">
  <Input
    bind:value={search}
    placeholder="Filter tasks..."
    class="h-8 w-[150px] lg:w-[250px]"
  />
  <p class="text-muted">Filters:</p>
  <Filters title="Severity" options={statuses} bind:selectedValues />
  <FiltersSingle title="Clients" options={statuses} bind:selectedValue />
  <Filters title="Tag" options={statuses} bind:selectedValues />
  <!-- Rule ID // maybe in search -->
  <!-- Time // maybe in search -->
</div>

<DataTable data={alerts} {columns} />

<div class="flex justify-between gap-4 flex-wrap">
  <p>
    Page {currentPage} of {maxPage}
  </p>
  <div>
    <Button
      variant="ghost"
      disabled={currentPage === 1}
      onclick={() => (currentPage = 1)}
    >
      First
    </Button>
    <Button
      variant="ghost"
      disabled={currentPage === 1}
      onclick={() => (currentPage -= 1)}
    >
      Previous
    </Button>

    <Button
      variant="ghost"
      disabled={currentPage === maxPage}
      onclick={() => (currentPage += 1)}
    >
      Next
    </Button>

    <Button
      variant="ghost"
      disabled={currentPage === maxPage}
      onclick={() => (currentPage = maxPage)}
    >
      Last
    </Button>
  </div>
</div>
