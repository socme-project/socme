<script lang="ts">
  import { page } from "$app/state";
  import { ShieldAlert } from "lucide-svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";

  let id = page.url.pathname.split("/").pop();

  let client: any = $state({});

  $effect(() => {
    axios
      .get("/api/clients/" + id, {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        client = res.data.client;
        console.log(client);
      })
      .catch(() => {
        toast.error("Internal server error");
      });
  });
</script>

<h1 class="flex items-center gap-4 mb-8">
  <ShieldAlert />
  Client {client.Name}
</h1>

<table class="w-full text-left">
  <tbody>
    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >id</th
      >
      <td>{client.ID}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >name</th
      >
      <td>{client.Name}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >wazuh host</th
      >
      <td>{client.WazuhIP}:{client.WazuhPort}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >wazuh version</th
      >
      <td>{client.WazuhVersion}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >wazuh username</th
      >
      <td>{client.WazuhUsername}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >wazuh password</th
      >
      <td>{client.WazuhPassword}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >indexer host</th
      >
      <td>{client.IndexerIP}:{client.IndexerPort}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >indexer username</th
      >
      <td>{client.IndexerUsername}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >indexer password</th
      >
      <td>{client.IndexerPassword}</td>
    </tr>

    <tr>
      <th
        class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
        >last alert</th
      >
      <td>{client.LastAlert}</td>
    </tr>
  </tbody>
</table>
