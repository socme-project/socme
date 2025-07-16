<script lang="ts">
  import Loading from "$src/lib/components/loading.svelte";
  import type { Client } from "$src/lib/stores/client";
  import { sendError } from "$src/lib/utils";
  import axios from "axios";
  import { route } from "sv-router/generated";
  import { onMount } from "svelte";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import { ExternalLink, Eye, EyeOff, Ghost } from "@lucide/svelte";
  import * as Card from "$lib/components/ui/card/index.js";
  import Clients from "$src/lib/components/charts/clients.svelte";
  const id = route.params.id;

  let client = $state<Client | null>(null);

  let showWazuhPassword = $state(false);
  let showIndexerPassword = $state(false);

  onMount(async () => {
    axios
      .get(`/api/clients/${id}`)
      .then((response) => {
        client = response.data.client;
      })
      .catch((error) => {
        sendError("Error fetching client:", error);
      });
  });
</script>

{#if client}
  <div class="flex flex-wrap gap-5 justify-start items-center my-4">
    <Avatar.Root class="h-14 w-14 rounded-full">
      <Avatar.Image
        class="rounded-full"
        src={client.Logo || ""}
        alt={client.Name + " avatar"}
      />
      <Avatar.Fallback class="rounded-full"><Ghost size={14} /></Avatar.Fallback
      >
    </Avatar.Root>
    <div>
      <div class="flex flex-wrap gap-2 items-center">
        <h1 class="m-0">{client.Name}</h1>
        <div class="relative flex ml-4 justify-center items-center">
          {#if client.WazuhIsAlive}
            <div class="absolute z-10 rounded-full w-2 h-2 bg-green-400"></div>
            <div
              class="absolute z-0 animate-pulse rounded-full w-4 h-4 bg-green-600/10"
            ></div>
          {:else}
            <div class="absolute z-10 rounded-full w-2 h-2 bg-red-500"></div>
            <div
              class="absolute z-0 animate-pulse rounded-full w-4 h-4 bg-red-500/20"
            ></div>
          {/if}
        </div>
      </div>
      <p class="text-muted-foreground text-xs">{client.ID}</p>
    </div>
  </div>

  <div class="grid xl:grid-cols-2 grid-cols-1 gap-5">
    <Card.Root class="flex flex-col">
      <Card.Header class="items-center">
        <Card.Title>Agents status</Card.Title>
      </Card.Header>
      <Card.Content class="flex-1 items-center flex">
        <Clients
          actif={client.ConnectedAgents || 0}
          inactif={client.DisconnectedAgents || 0}
        />
      </Card.Content>
    </Card.Root>

    <Card.Root class="flex flex-col">
      <Card.Header class="items-center">
        <Card.Title>Client informations</Card.Title>
      </Card.Header>
      <Card.Content class="flex-1">
        <table class="w-full text-left">
          <tbody>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Wazuh version</th
              >
              <td>{client.WazuhVersion || "version unknown"}</td>
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Last alert</th
              >
              <td
                >{client.LastAlert === "0001-01-01T00:00:00Z"
                  ? "timestamp unknown"
                  : client.LastAlert}</td
              >
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Host</th
              >
              <td>{client.Host}</td>
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Wazuh</th
              >
              <td
                ><a
                  class="flex gap-2 items-center hover:underline"
                  href={client.Host + ":" + client.WazuhPort}
                  >{client.Host}:{client.WazuhPort}
                  <ExternalLink size={12} /></a
                ></td
              >
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Wazuh username</th
              >
              <td>{client.WazuhUsername}</td>
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Wazuh password</th
              >
              <td
                >{#if showWazuhPassword}
                  <p class="flex gap-2 items-center">
                    {client.WazuhPassword}<Eye
                      onclick={() => (showWazuhPassword = !showWazuhPassword)}
                      size={12}
                    />
                  </p>
                {:else}
                  <p class="flex gap-2 items-center text-muted-foreground">
                    Password redacted
                    <EyeOff
                      onclick={() => (showWazuhPassword = !showWazuhPassword)}
                      size={12}
                    />
                  </p>
                {/if}
              </td>
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Indexer</th
              >
              <td
                ><a
                  class="flex gap-2 items-center hover:underline"
                  href={client.Host + ":" + client.IndexerPort}
                  >{client.Host}:{client.IndexerPort}
                  <ExternalLink size={12} /></a
                ></td
              >
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Indexer username</th
              >
              <td>{client.IndexerUsername}</td>
            </tr>
            <tr>
              <th
                class="w-[220px] text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
                >Indexer password</th
              >
              <td
                >{#if showIndexerPassword}
                  <p class="flex gap-2 items-center">
                    {client.IndexerPassword}<Eye
                      onclick={() =>
                        (showIndexerPassword = !showIndexerPassword)}
                      size={12}
                    />
                  </p>
                {:else}
                  <p class="flex gap-2 items-center text-muted-foreground">
                    Password redacted
                    <EyeOff
                      onclick={() =>
                        (showIndexerPassword = !showIndexerPassword)}
                      size={12}
                    />
                  </p>
                {/if}
              </td>
            </tr>
          </tbody>
        </table>
      </Card.Content>
    </Card.Root>
  </div>
{:else}
  <Loading />
{/if}
