<script lang="ts">
  import {
    ChevronLeft,
    ChevronRight,
    ChevronsLeft,
    ChevronsRight,
    FileQuestionIcon,
    ShieldAlert,
    TouchpadOffIcon,
  } from "lucide-svelte";
  import { columns } from "./columns";
  import DataTable from "./data-table.svelte";
  import type { Alert } from "./columns";
  import Filters from "./filters.svelte";
  import FiltersSingle from "./filtersSingle.svelte";
  import Input from "$lib/components/ui/input/input.svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Select from "$lib/components/ui/select/index.js";

  let currentPage = $state(1);

  let perPage = $state(10);
  let perPageString = $state("10");
  $effect(() => {
    perPage = parseInt(perPageString);
  });
  let maxPage = $state(1);

  let alerts: Alert[] = $state([]);

  $effect(() => {
    axios
      .get("/api/alerts", {
        headers: { Authorization: localStorage.getItem("token") },
        params: {
          page: currentPage,
          perPage: perPage,
          severity: severityValues.join(","),
          search: search,
        },
      })
      .then((res) => {
        console.log(res.data);
        alerts = res.data.alerts.map((alert: any) => ({
          id: alert.ID,
          title: alert.rule_description,
          severity: alert.rule_level,
          client: alert.client_name,
          timestamp: alert.timestamp,
          raw: alert.raw_json,
        }));
        maxPage = res.data.maxPage;
      })
      .catch(() => {
        toast.error("Internal server error");
      });
  });

  export const severity = [
    {
      value: "low",
      label: "Low",
      icon: FileQuestionIcon,
    },
    {
      value: "medium",
      label: "Medium",
      icon: TouchpadOffIcon,
    },
    {
      value: "high",
      label: "High",
      icon: FileQuestionIcon,
    },
    {
      value: "critical",
      label: "Critical",
      icon: TouchpadOffIcon,
    },
  ];

  let selectedValues = $state([]);
  let severityValues = $state([]);
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
  <Filters
    title="Severity"
    options={severity}
    bind:selectedValues={severityValues}
  />
  <FiltersSingle title="Clients" options={severity} bind:selectedValue />
  <Filters title="Tag" options={severity} bind:selectedValues />
  <!-- Rule ID // maybe in search -->
  <!-- Time // maybe in search -->
</div>

<DataTable data={alerts} {columns} />

<div class="flex justify-between gap-4 flex-wrap items-center my-4">
  <div class="flex items-center gap-4">
    <p class="text-xs text-muted-foreground">
      Page {currentPage} of {maxPage} | Per page:
    </p>
    <Select.Root type="single" bind:value={perPageString}>
      <Select.Trigger class="w-[80px] p-1 h-6">{perPage}</Select.Trigger>
      <Select.Content>
        <Select.Item value="10">10</Select.Item>
        <Select.Item value="20">20</Select.Item>
        <Select.Item value="50">50</Select.Item>
        <Select.Item value="100">100</Select.Item>
      </Select.Content>
    </Select.Root>
  </div>
  <div>
    <Button
      variant="ghost"
      disabled={currentPage === 1}
      onclick={() => (currentPage = 1)}
    >
      <ChevronsLeft />
    </Button>
    <Button
      variant="ghost"
      disabled={currentPage === 1}
      onclick={() => (currentPage -= 1)}
    >
      <ChevronLeft />
    </Button>

    <Button
      variant="ghost"
      disabled={currentPage === maxPage}
      onclick={() => (currentPage += 1)}
    >
      <ChevronRight />
    </Button>

    <Button
      variant="ghost"
      disabled={currentPage === maxPage}
      onclick={() => (currentPage = maxPage)}
    >
      <ChevronsRight />
    </Button>
  </div>
</div>
