<script lang="ts">
  import axios from "axios";
  import { sendError } from "../utils";
  import { toast } from "svelte-sonner";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import Button from "./ui/button/button.svelte";
  import type { Client } from "$src/lib/stores/client";
  const { client, runAfterCreation } = $props();

  function editClient(client: Client) {
    axios
      .patch(`/api/clients/${client.ID}`, null, {
        params: {
          name: client.Name,
          logo: client.Logo,
          wazuh_ip: client.WazuhIP,
          wazuh_port: client.WazuhPort,
          wazuh_username: client.WazuhUsername,
          wazuh_password: client.WazuhPassword,
          indexer_ip: client.IndexerIP,
          indexer_port: client.IndexerPort,
          indexer_username: client.IndexerUsername,
          indexer_password: client.IndexerPassword,
        },
      })
      .then(() => {
        toast.success("Client edited successfully.");
        if (!runAfterCreation) return;
        runAfterCreation();
      })
      .catch((error) => sendError("Error editing client:", error));
  }
</script>

<div class="grid gap-4 py-4">
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="name" class="text-right">Name</Label>
    <Input id="name" bind:value={client.Name} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="logo" class="text-right">Logo URL</Label>
    <Input id="logo" bind:value={client.Logo} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_ip" class="text-right">Wazuh IP</Label>
    <Input id="wazuh_ip" bind:value={client.WazuhIP} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_port" class="text-right">Wazuh Port</Label>
    <Input id="wazuh_port" bind:value={client.WazuhPort} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_username" class="text-right">Wazuh Username</Label>
    <Input
      id="wazuh_username"
      bind:value={client.WazuhUsername}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_password" class="text-right">Wazuh Password</Label>
    <Input
      type="password"
      id="wazuh_password"
      bind:value={client.WazuhPassword}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_ip" class="text-right">Indexer IP</Label>
    <Input id="indexer_ip" bind:value={client.IndexerIP} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_port" class="text-right">Indexer Port</Label>
    <Input
      id="indexer_port"
      bind:value={client.IndexerPort}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_username" class="text-right">Indexer Username</Label>
    <Input
      id="indexer_username"
      bind:value={client.IndexerUsername}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_password" class="text-right">Indexer Password</Label>
    <Input
      type="password"
      id="indexer_password"
      bind:value={client.IndexerPassword}
      class="col-span-3"
    />
  </div>
</div>

<Button onclick={() => editClient(client)} class="w-full mt-6"
  >Edit client</Button
>
