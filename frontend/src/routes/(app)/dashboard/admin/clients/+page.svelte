<script lang="ts">
  import { Crown, FolderCog } from "lucide-svelte";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import Ellipsis from "lucide-svelte/icons/ellipsis";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import { onMount } from "svelte";
  import NewClient from "./newClient.svelte";
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";

  let clients: any = $state([]);

  onMount(async () => {
    await axios
      .get("/api/clients/list", {
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
  <NewClient />
</div>

<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-[300px]">name</Table.Head>
      <Table.Head>artemis is alive</Table.Head>
      <Table.Head>artemis version</Table.Head>
      <Table.Head>artemis IP</Table.Head>
      <Table.Head></Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#if clients.length === 0}
      {#each Array(6) as _}
        <Table.Row>
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
          <Table.Cell class="flex justify-start items-center gap-4">
            <div
              class="rounded-full animate-pulse h-2 w-2"
              style={"background-color: " +
                (client.ArtemisIsAlive ? "#4ADE80;" : "#F87171;")}
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
          <Table.Cell>{client.ArtemisIsAlive}</Table.Cell>
          <Table.Cell>{client.ArtemisVersion}</Table.Cell>
          <Table.Cell>{client.ArtemisIP}</Table.Cell>

          <Table.Cell class="text-right w-10">
            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                <Ellipsis size={18} />
              </DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Group>
                  <!-- TODO: Enable buttons -->
                  <DropdownMenu.Item>More informations</DropdownMenu.Item>
                  <DropdownMenu.Item>Open Wazuh dashboard</DropdownMenu.Item>
                  <DropdownMenu.Item>Upgrade artemis</DropdownMenu.Item>
                  <DropdownMenu.Item>See artemis password</DropdownMenu.Item>
                </DropdownMenu.Group>
                <DropdownMenu.Separator />
                <DropdownMenu.Group>
                  <DropdownMenu.Item>Delete client</DropdownMenu.Item>
                </DropdownMenu.Group>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </Table.Cell>
        </Table.Row>
      {/each}
    {/if}
  </Table.Body>
</Table.Root>
