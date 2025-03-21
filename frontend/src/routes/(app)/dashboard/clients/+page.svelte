<script lang="ts">
  import { Crown, FolderCog } from "lucide-svelte";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import Ellipsis from "lucide-svelte/icons/ellipsis";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import { onMount } from "svelte";
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";

  let clients: any = $state([]);

  onMount(async () => {
    await axios
      .get("/api/clients", {
        headers: { Authorization: localStorage.getItem("token") },
      })
      .then((res) => {
        clients = res.data.clients;
      })
      .catch(() => {
        toast.error("Internal server error");
      });
  });
</script>

<div class="w-full flex justify-between flex-wrap items-center">
  <h1 class="flex items-center gap-4">
    <FolderCog />
    Clients
  </h1>
</div>

<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-[30px]">ID</Table.Head>
      <Table.Head class="w-[300px]">name</Table.Head>
      <Table.Head>version</Table.Head>
      <Table.Head>wazuh host</Table.Head>
      <Table.Head>indexer host</Table.Head>
      <Table.Head></Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#if clients.length === 0}
      {#each Array(6) as _}
        <Table.Row>
          <Table.Cell>
            <Skeleton class="w-[30px] h-4" />
          </Table.Cell>
          <Table.Cell class="flex justify-start items-center gap-4">
            <Skeleton class="rounded-full animate-pulse h-2 w-2" />
            <Skeleton class="size-10 rounded-full" />
            <Skeleton class="w-[100px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[120px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[70px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[90px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[20px] h-3" />
          </Table.Cell>
        </Table.Row>
      {/each}
    {:else}
      {#each clients as client}
        <Table.Row>
          <Table.Cell>{client.ID}</Table.Cell>
          <Table.Cell class="flex justify-start items-center gap-4">
            <div
              class="rounded-full animate-pulse h-2 w-2"
              style={"background-color: " +
                (client.WazuhIsAlive ? "#4ADE80;" : "#F87171;")}
            ></div>
            <Avatar.Root>
              <Avatar.Image
                class="object-cover"
                src={client.Logo}
                alt={client.Name + " logo"}
              />
              <Avatar.Fallback>
                <Ellipsis />
              </Avatar.Fallback>
            </Avatar.Root>
            <div class="flex items-center gap-2">
              {client.Name}
            </div></Table.Cell
          >
          <Table.Cell>{client.Version}</Table.Cell>
          <Table.Cell>{client.WazuhIP + ":" + client.WazuhPort}</Table.Cell>
          <Table.Cell>{client.IndexerIP + ":" + client.IndexerPort}</Table.Cell>

          <Table.Cell class="text-right w-10">
            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                <Ellipsis size={18} />
              </DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Group>
                  <a href={"/dashboard/clients/" + client.ID}>
                    <DropdownMenu.Item>More informations</DropdownMenu.Item>
                  </a>
                  <a href={"https://" + client.WazuhIP + ":" + 443}>
                    <DropdownMenu.Item>Open Wazuh dashboard</DropdownMenu.Item>
                  </a>
                </DropdownMenu.Group>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </Table.Cell>
        </Table.Row>
      {/each}
    {/if}
  </Table.Body>
</Table.Root>
