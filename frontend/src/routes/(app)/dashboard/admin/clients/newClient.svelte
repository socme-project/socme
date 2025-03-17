<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import axios from "axios";
  import { toast } from "svelte-sonner";

  let name = $state("");
  let logo = $state("");

  let wazuhHost = $state(""); // 192.168.1.102:55000
  let wazuhUsername = $state("");
  let wazuhPassword = $state("");

  let indexerHost = $state("");
  let indexerUsername = $state("");
  let indexerPassword = $state("");

  let isOpen = $state(false);

  async function handleSubmit() {
    if (
      !name ||
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
      .get("/api/client/new", {
        headers: { Authorization: localStorage.getItem("token") },
        params: {
          name: name,
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
      })
      .then(() => {
        toast.success("Client created");
        isOpen = false;
      })
      .catch((error) => {
        toast.error("Failed to create client");
        console.log(error);
      });
  }
  // TODO: Ask if we want to generate a new password
</script>

<Dialog.Root bind:open={isOpen}>
  <Dialog.Trigger class={buttonVariants({ variant: "outline" })}
    >New client</Dialog.Trigger
  >
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>New client</Dialog.Title>
      <Dialog.Description>
        Create a new client by filling out the form below.
      </Dialog.Description>
    </Dialog.Header>
    <div class="grid gap-4 py-4">
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="name" class="text-right">Name</Label>
        <Input
          id="name"
          bind:value={name}
          class="col-span-3"
          placeholder="Client name"
        />
      </div>
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
        <label for="indexerusername" class="text-right">Indexer username</label>
        <Input
          id="indexerusername"
          bind:value={indexerUsername}
          class="col-span-3"
          placeholder="admin"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <label for="indexerpassword" class="text-right">Indexer password</label>
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
      <Button onclick={handleSubmit}>Create the client</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
