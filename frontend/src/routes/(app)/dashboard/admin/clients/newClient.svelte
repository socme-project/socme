<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  let name = $state("");
  let logo = $state("");
  let artemisIP = $state("");
  let artemisPassword = $state("");

  let isOpen = $state(false);

  async function handleSubmit() {
    await axios
      .get("/api/client/new", {
        headers: { Authorization: localStorage.getItem("token") },
        params: {
          name: name,
          logo: logo,
          artemisIP: artemisIP,
          artemisPassword: artemisPassword,
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
        <Label for="artemisIP" class="text-right">Artemis IP</Label>
        <Input
          id="artemisIP"
          bind:value={artemisIP}
          class="col-span-3"
          placeholder="Artemis IP"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="artemisPassword" class="text-right">Artemis password</Label>
        <Input
          id="artemisPassword"
          bind:value={artemisPassword}
          class="col-span-3"
          placeholder="Artemis password"
        />
      </div>
    </div>
    <Dialog.Footer>
      <Button onclick={handleSubmit}>Create the client</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
