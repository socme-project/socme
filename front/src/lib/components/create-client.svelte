<script lang="ts">
  import axios from "axios";
  import { sendError } from "../utils";
  import { toast } from "svelte-sonner";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import Button from "./ui/button/button.svelte";
  import type { Client } from "$src/lib/stores/client";

  const client = $state<Client>({
    ID: "",
    Name: "",
    Logo: "",

    Host: "",

    SshPort: "",
    SshUsername: "",
    SshPassword: "",

    WazuhPort: "",
    WazuhUsername: "",
    WazuhPassword: "",

    IndexerPort: "",
    IndexerUsername: "",
    IndexerPassword: "",

    LastAlert: "",
    WazuhIsAlive: false,
    WazuhVersion: "",
  });

  function createClient(client: Client) {
    axios
      .post("/api/clients", null, {
        params: {
          name: client.Name,
          logo: client.Logo,

          host: client.Host,

          ssh_port: client.SshPort,
          ssh_username: client.SshUsername,
          ssh_password: client.SshPassword,

          wazuh_port: client.WazuhPort,
          wazuh_username: client.WazuhUsername,
          wazuh_password: client.WazuhPassword,

          indexer_port: client.IndexerPort,
          indexer_username: client.IndexerUsername,
          indexer_password: client.IndexerPassword,
        },
      })
      .then(() => {
        toast.success("Client created successfully.");
        if (!runAfterCreation) return;
        runAfterCreation();
      })
      .catch((error) => sendError("Error creating client:", error));
  }

  const { runAfterCreation } = $props();
</script>

<div class="grid gap-4 py-4">
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="name">Name</Label>
    <Input id="name" bind:value={client.Name} placeholder="acme inc" class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="logo">Logo URL</Label>
    <Input id="logo" bind:value={client.Logo} placeholder="https://example.com/mylogo.png" class="col-span-3" />
  </div>

  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="host">Host</Label>
    <Input id="host" bind:value={client.Host} placeholder="Ip or domain" class="col-span-3" />
  </div>

  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="ssh_port">Ssh Port</Label>
    <Input id="ssh_port" bind:value={client.SshPort} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="ssh_username">Ssh Username</Label>
    <Input
      id="ssh_username"
      bind:value={client.SshUsername}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="ssh_password">Ssh Password</Label>
    <Input
      type="password"
      id="ssh_password"
      bind:value={client.SshPassword}
      class="col-span-3"
    />
  </div>

  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_port">Wazuh Port</Label>
    <Input id="wazuh_port" bind:value={client.WazuhPort} class="col-span-3" />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_username">Wazuh Username</Label>
    <Input
      id="wazuh_username"
      bind:value={client.WazuhUsername}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="wazuh_password">Wazuh Password</Label>
    <Input
      type="password"
      id="wazuh_password"
      bind:value={client.WazuhPassword}
      class="col-span-3"
    />
  </div>

  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_port">Indexer Port</Label>
    <Input
      id="indexer_port"
      bind:value={client.IndexerPort}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_username">Indexer Username</Label>
    <Input
      id="indexer_username"
      bind:value={client.IndexerUsername}
      class="col-span-3"
    />
  </div>
  <div class="grid grid-cols-4 items-center gap-4">
    <Label for="indexer_password">Indexer Password</Label>
    <Input
      type="password"
      id="indexer_password"
      bind:value={client.IndexerPassword}
      class="col-span-3"
    />
  </div>
</div>

<Button onclick={() => createClient(client)} class="w-full mt-6"
  >Create client</Button
>
