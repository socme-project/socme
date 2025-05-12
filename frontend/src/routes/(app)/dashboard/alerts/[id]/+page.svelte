<script lang="ts">
  import { page } from "$app/state";
  import { ShieldAlert } from "lucide-svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";

  let id = page.url.pathname.split("/").pop();

  let alert: any = $state({});

  $effect(() => {
    axios
      .get("/api/alerts/" + id, {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        alert = res.data.alert;
        console.log(alert);
      })
      .catch(() => {
        toast.error("Internal server error");
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

<h1 class="flex items-center gap-4 mb-8">
  <ShieldAlert />
  Alert #{id}
</h1>

<h2>Alert details</h2>

<table class="w-full text-left">
  <tbody>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >id</th
      >
      <td>{alert.ID}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Severity</th
      >
      <td>{alert.rule_level}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Description</th
      >
      <td>{alert.rule_description}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Client</th
      >
      <td>{alert.client_name}</td>
    </tr>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >Time</th
      >
      <td>{alert.timestamp}</td>
    </tr>
  </tbody>
</table>

<h2>Raw data</h2>

<table class="w-full text-left">
  <tbody>
    {#if alert.raw_json !== undefined}
      {#each Object.entries(flattenObject(JSON.parse(alert.raw_json))) as [key, value]}
        <tr>
          <th
            class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
            >{key}</th
          >
          <td>{value}</td>
        </tr>
      {/each}
    {/if}
  </tbody>
</table>
