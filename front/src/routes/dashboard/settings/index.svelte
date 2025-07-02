<script lang="ts">
  import Button from "$src/lib/components/ui/button/button.svelte";
  import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import Input from "$lib/components/ui/input/input.svelte";
  import { setMode, resetMode } from "mode-watcher";
  import Label from "$lib/components/ui/label/label.svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import { sendError } from "$src/lib/utils";
  import { user } from "$src/lib/stores/user";
  import * as Tabs from "$lib/components/ui/tabs/index.js";

  import { PencilOff, Settings } from "@lucide/svelte";
    import { navigate } from "sv-router/generated";

  function deleteUser(id: string) {
    axios
      .delete(`/api/users/${id}`)
      .then(() => {
        toast.success("User deleted successfully.");
        navigate("/auth/logout")
      })
      .catch((error) => sendError("Error deleting user:", error));
  }
</script>

<h1 class="flex items-center gap-4">
  <Settings />
  Settings
</h1>

<Tabs.Root value="account" class="w-full">
  <Tabs.List
  >
    <Tabs.Trigger value="account"
      >Account</Tabs.Trigger
    >
    <Tabs.Trigger value="appearance"
      >Appearance</Tabs.Trigger
    >
  </Tabs.List>
  <Tabs.Content value="account">
    Make changes to your account here.
    <h2 class="text-xl font-semibold">Account informations</h2>
    {#if $user}
      <div class="flex flex-wrap py-4 gap-12">
        <div>
          <div class="relative">
            <Avatar.Root
              class="h-20 w-20 md:h-40 md:w-40 rounded-full p-1 border-muted-foreground border"
            >
              <Avatar.Image
                class="rounded-full"
                src={"https://avatars.githubusercontent.com/u/" +
                  $user.GithubID +
                  "?v=4"}
                alt={$user.Name + " avatar"}
              />
              <Avatar.Fallback class="rounded-full"
                >{$user.Name}</Avatar.Fallback
              >
            </Avatar.Root>
            <PencilOff
              class="absolute bottom-0 right-[-8px] text-muted-foreground"
            />
          </div>
        </div>
        <div class="flex flex-col gap-4">
          <div class="flex w-full max-w-sm flex-col gap-1.5">
            <Label for="username">Username</Label>
            <Input type="text" id="username" value={$user.Name} disabled />
          </div>
          <div class="flex w-full max-w-sm flex-col gap-1.5">
            <Label for="githubid">Github ID</Label>
            <Input type="text" id="githubid" value={$user.GithubID} disabled />
          </div>
          <div class="flex w-full max-w-sm flex-col gap-1.5">
            <Label for="role">Role</Label>
            <Input type="text" id="role" value={$user.Role} disabled />
          </div>
        </div>
      </div>
      <p class="text-muted-foreground text-xs pt-3">
        (!) Account information are in read-only mode. Change information from
        the Auth provider (github).
      </p>
    {/if}
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <Button variant="destructive">Delete account</Button>
      </AlertDialog.Trigger>
      <AlertDialog.Content>
        <AlertDialog.Header>
          <AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
          <AlertDialog.Description>
            This action cannot be undone. This will permanently delete the
            account and remove the data from the servers.
          </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
          <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
          {#if $user}
            <AlertDialog.Action
              onclick={() => {
                deleteUser($user.ID);
              }}>Continue</AlertDialog.Action
            >
          {/if}
        </AlertDialog.Footer>
      </AlertDialog.Content>
    </AlertDialog.Root>
  </Tabs.Content>
  <Tabs.Content value="appearance">
    Change the appearance of your account here.

    <h2 class="text-xl font-semibold">Theme</h2>

    <div class="flex flex-wrap gap-5 max-w-3xl">
      <Label onclick={() => setMode("light")}>
        <div
          class="border-primary dark:border-muted hover:border-accent items-center rounded-md border-2 p-1"
        >
          <div class="space-y-2 rounded-sm bg-[#ecedef] p-2">
            <div class="space-y-2 rounded-md bg-white p-2 shadow-sm">
              <div class="h-2 w-[80px] rounded-lg bg-[#ecedef]"></div>
              <div class="h-2 w-[100px] rounded-lg bg-[#ecedef]"></div>
            </div>
            <div
              class="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm"
            >
              <div class="h-4 w-4 rounded-full bg-[#ecedef]"></div>
              <div class="h-2 w-[100px] rounded-lg bg-[#ecedef]"></div>
            </div>
            <div
              class="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm"
            >
              <div class="h-4 w-4 rounded-full bg-[#ecedef]"></div>
              <div class="h-2 w-[100px] rounded-lg bg-[#ecedef]"></div>
            </div>
          </div>
        </div>
      </Label>
      <Label onclick={() => setMode("dark")}>
        <div
          class="dark:border-primary border-muted bg-popover hover:bg-accent hover:text-accent-foreground items-center rounded-md border-2 p-1"
        >
          <div class="space-y-2 rounded-sm bg-slate-950 p-2">
            <div class="space-y-2 rounded-md bg-slate-800 p-2 shadow-sm">
              <div class="h-2 w-[80px] rounded-lg bg-slate-400"></div>
              <div class="h-2 w-[100px] rounded-lg bg-slate-400"></div>
            </div>
            <div
              class="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm"
            >
              <div class="h-4 w-4 rounded-full bg-slate-400"></div>
              <div class="h-2 w-[100px] rounded-lg bg-slate-400"></div>
            </div>
            <div
              class="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm"
            >
              <div class="h-4 w-4 rounded-full bg-slate-400"></div>
              <div class="h-2 w-[100px] rounded-lg bg-slate-400"></div>
            </div>
          </div>
        </div>
      </Label>
    </div>
    <Button
      class="text-muted-foreground"
      variant="link"
      onclick={() => resetMode()}>or system theme</Button
    >
  </Tabs.Content>
</Tabs.Root>
