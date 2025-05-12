<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import axios from "axios";
  import { toast } from "svelte-sonner";

  let {
    clientId,

    logo = "",
    wazuhHost = "",
    wazuhUsername = "",
    wazuhPassword = "",
    indexerHost = "",
    indexerUsername = "",
    indexerPassword = "",
  } = $props();

  let isOpen = $state(false);

  async function handleSubmit() {
    if (
      !logo ||
      !wazuhHost ||
      !wazuhUsername ||
      !wazuhPassword ||
      !indexerHost ||
      !indexerUsername ||
      !indexerPassword
    ) {
      toast.error("Please fill out all fields");
      return;
    }
    if (!wazuhHost.includes(":") || !indexerHost.includes(":")) {
      toast.error("Please provide a valid host:port combination");
      return;
    }

    await axios
      .patch(
        "/api/client/" + clientId,
        {
          clientId: clientId,
          logo: logo,

          wazuhIP: wazuhHost.split(":")[0],
          wazuhPort: wazuhHost.split(":")[1],
          wazuhUsername: wazuhUsername,
          wazuhPassword: wazuhPassword,

          indexerIP: indexerHost.split(":")[0],
          indexerPort: indexerHost.split(":")[1],
          indexerUsername: indexerUsername,
          indexerPassword: indexerPassword,
        },
        {
          headers: {
            Authorization: localStorage.getItem("token"),
            "Content-Type": "application/json",
          },
        },
      )
      .then(() => {
        toast.success("Client edited");
        isOpen = false;
      })
      .catch((error) => {
        const errorMessage =
          error.response?.data?.error || "An unknown error occurred";
        toast.error(`Failed to edit client: ${errorMessage}`);
        console.log(error);
      });
  }
</script>

<Dialog.Root bind:open={isOpen}>
  <Dialog.Trigger class={buttonVariants({ variant: "outline" })}
    >Edit client</Dialog.Trigger
  >
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>Edit client</Dialog.Title>
      <Dialog.Description>Edit the client information</Dialog.Description>
    </Dialog.Header>
    <div class="grid gap-4 py-4">
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="logo" class="text-right">Logo</Label>
        <Input
          id="logo"
          bind:value={logo}
          class="col-span-3"
          placeholder="Client logo"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="wazuhhost" class="text-right">Wazuh host</Label>
        <Input
          id="wazuhhost"
          bind:value={wazuhHost}
          class="col-span-3"
          placeholder="192.168.1.102:55000"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="wazuhusername" class="text-right">Wazuh username</Label>
        <Input
          id="wazuhusername"
          bind:value={wazuhUsername}
          class="col-span-3"
          placeholder="admin"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="wazuhpassword" class="text-right">Wazuh password</Label>
        <Input
          id="wazuhpassword"
          bind:value={wazuhPassword}
          class="col-span-3"
          placeholder="mypassword123"
          type="password"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="indexerhost" class="text-right">Indexer host</Label>
        <Input
          id="indexerhost"
          bind:value={indexerHost}
          class="col-span-3"
          placeholder="192.168.1.102:9200"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="indexerusername" class="text-right">Indexer username</Label>
        <Input
          id="indexerusername"
          bind:value={indexerUsername}
          class="col-span-3"
          placeholder="admin"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="indexerpassword" class="text-right">Indexer password</Label>
        <Input
          id="indexerpassword"
          bind:value={indexerPassword}
          type="password"
          class="col-span-3"
          placeholder="mypassword123"
        />
      </div>
    </div>
    <Dialog.Footer>
      <Button onclick={handleSubmit}>Save</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
