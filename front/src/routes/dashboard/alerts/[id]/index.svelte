<script lang="ts">
  import { sendError } from "$src/lib/utils";
  import axios from "axios";
  import { onMount } from "svelte";
  import { route } from "sv-router/generated";

  import Loading from "$src/lib/components/loading.svelte";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import { ExternalLink, Eye, EyeOff, Ghost } from "@lucide/svelte";
  import type { Alert } from "$src/lib/components/alerts/columns";

  const id = route.params.id;

  let alert = $state<Alert|null>(null);

  onMount(async () => {
    axios
      .get(`/api/alerts/${id}`)
      .then((response) => {
        alert = response.data.alert;
        console.log(alert)
      })
      .catch((error) => {
        sendError("Error fetching alert:", error);
      });
  });

  function flattenObject(obj: any, prefix = "") {
    let result: Record<string, any> = {};
    for (const [key, value] of Object.entries(obj)) {
      const newKey = prefix ? `${prefix}.${key}` : key;
      if (typeof value === "object" && value !== null) {
        Object.assign(result, flattenObject(value, newKey));
      } else {
        result[newKey] = value;
      }
    }
    return result;
  }
</script>

{#if alert}
  <div class="flex flex-wrap gap-5 justify-start items-center mt-4">
    <Avatar.Root class="h-14 w-14 rounded-full">
      <Avatar.Image
        class="rounded-full"
        src={alert.Client.Logo || ""}
        alt={alert.Client.Name + " avatar"}
      />
      <Avatar.Fallback class="rounded-full"><Ghost size={14} /></Avatar.Fallback
      >
    </Avatar.Root>
    <div>
      <div class="flex flex-wrap gap-2 items-center">
        <h1 class="m-0">{alert.ID}</h1>
      </div>
      <p class="text-muted-foreground text-xs">{alert.Client.ID}</p>
    </div>
  </div>

  <table class="w-full text-left mt-10">
    <tbody>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >Client</th
        >
        <td class="flex gap-2 items-center h-12">
            <Avatar.Root class="h-6 w-6 rounded-full mr-4">
              <Avatar.Image
                class="rounded-full"
                src={alert.Client.Logo || ""
                  }
                alt={alert.Client.Name + " avatar"}
              />
              <Avatar.Fallback class="rounded-full"><Ghost size={14}/></Avatar.Fallback
              >
            </Avatar.Root>

            <div class="relative  h-full flex mr-2 justify-center items-center">
              {#if alert.Client.WazuhIsAlive}
                <div
                  class="absolute z-10 rounded-full w-2 h-2 bg-green-400"
                ></div>
                <div
                  class="absolute z-0 animate-pulse rounded-full w-4 h-4 bg-green-600/10"
                ></div>
              {:else}
                <div
                  class="absolute z-10 rounded-full w-2 h-2 bg-red-500"
                ></div>
                <div
                  class="absolute z-0 animate-pulse rounded-full w-4 h-4 bg-red-500/20"
                ></div>
              {/if}
            </div>
            <a href={`/dashboard/clients/${alert.Client.ID}`}>
            {alert.Client.Name}
            </a>
        </td>
      </tr>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >ID</th
        >
        <td>{alert.ID}</td>
      </tr>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >Wazuh alert ID</th
        >
        <td>{alert.WazuhAlertID}</td>
      </tr>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >Description</th
        >
        <td>{alert.RuleDescription}</td>
      </tr>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >Rule ID</th
        >
        <td>{alert.RuleID}</td>
      </tr>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >Rule level</th
        >
        <td>{alert.RuleLevel}</td>
      </tr>
      <tr>
        <th
          class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
          >Timestamp</th
        >
        <td>{alert.Timestamp}</td>
      </tr>
    </tbody>
  </table>

  <h2>Raw</h2>

<table class="w-full text-left">
  <tbody>
      {#each Object.entries(flattenObject(JSON.parse(alert.RawJSON))) as [key, value]}
        <tr>
          <th
            class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
            >{key}</th
          >
          <td>{value}</td>
        </tr>
      {/each}
  </tbody>
</table>
{:else}
  <Loading />
{/if}
