<script lang="ts">
  import { Ghost, Folder } from "@lucide/svelte";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import axios from "axios";
  import { onMount } from "svelte";
  import * as Table from "$lib/components/ui/table/index.js";
  import { sendError } from "$src/lib/utils";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import type { Client } from "$src/lib/stores/client";
  import { navigate } from "sv-router/generated";

  interface ClientWithOpen extends Client {
    Open: boolean;
  }

  let clients = $state<ClientWithOpen[] | null>(null);

  function fetchClient() {
    axios
      .get("/api/clients")
      .then((response) => {
        clients = response.data.clients.map((client: Client) => ({
          ...client,
          Open: false, // Initialize dropdown state to closed
        }));
      })
      .catch((error) => sendError("Error fetching clients:", error));
  }

  onMount(async () => {
    fetchClient();
  });
</script>

<div class="flex justify-between flex-wrap gap-8 items-center">
  <div>
    <h1 class="flex items-center gap-4">
      <Folder />
      Clients
    </h1>

    <p class="mb-4 text-muted-foreground">
      Manage clients here. You can view their details, and manage their configurations.
    </p>
  </div>
</div>

<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-[300px]">name</Table.Head>
      <Table.Head>wazuh</Table.Head>
      <Table.Head>indexer</Table.Head>
      <Table.Head>version</Table.Head>
      <Table.Head>actions</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#if clients === null}
      <Table.Row>
        <Table.Cell
          colspan={4}
          class="text-center animate-pulse text-muted-foreground"
        >
          Loading clients...
        </Table.Cell>
      </Table.Row>
    {:else if clients.length === 0}
      <Table.Row>
        <Table.Cell colspan={4} class="text-center text-muted-foreground"
          >No clients found.</Table.Cell
        >
      </Table.Row>
    {:else}
      {#each clients as client (client.ID)}
        <Table.Row>
          <Table.Cell class="flex gap-2 items-center">
            <Avatar.Root class="h-6 w-6 rounded-full mr-4">
              <Avatar.Image
                class="rounded-full"
                src={client.Logo || ""
                  }
                alt={client.Name + " avatar"}
              />
              <Avatar.Fallback class="rounded-full"><Ghost size={14}/></Avatar.Fallback
              >
            </Avatar.Root>

            <div class="relative flex mr-2 justify-center items-center">
              {#if client.WazuhIsAlive}
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
            <a href={`/dashboard/clients/${client.ID}`}>
            {client.Name}
            </a>
          </Table.Cell
          >
          <Table.Cell>
            {client.Host}:{client.WazuhPort}</Table.Cell
          >
          <Table.Cell>
            {client.Host}:{client.IndexerPort}</Table.Cell
          >
          <Table.Cell>
            {client.WazuhVersion}</Table.Cell
          >

          <Table.Cell>
            <DropdownMenu.Root bind:open={client.Open}>
              <DropdownMenu.Trigger>...</DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Group>
                  <DropdownMenu.Label>{client.Name}</DropdownMenu.Label>
                  <DropdownMenu.Separator />
                  <DropdownMenu.Item
                    onclick={() => navigate(`/dashboard/clients/${client.ID}`)}
                    >Details</DropdownMenu.Item
                  >
                </DropdownMenu.Group>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </Table.Cell>
        </Table.Row>
      {/each}
    {/if}
  </Table.Body>
</Table.Root>
