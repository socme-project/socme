<script lang="ts">
  import { UserCog } from "@lucide/svelte";
  import axios from "axios";
  import { onMount } from "svelte";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import { type User } from "$src/lib/stores/user";
  import { toast } from "svelte-sonner";
  import { sendError } from "$src/lib/utils";

  interface UserWithOpen extends User {
    Open: boolean;
  }

  let users = $state<UserWithOpen[] | null>(null);

  onMount(async () => {
    axios
      .get("/api/users")
      .then((response) => {
        users = response.data.users.map((user: User) => ({
          ...user,
          Open: false, // Initialize dropdown state to closed
        }));
      })
      .catch((error) => sendError("Error fetching users:", error));
  });

  function changeRole(id: string, newRole: string) {
    axios
      .patch(`/api/users/${id}/role`, null, { params: { role: newRole } })
      .then(() => {
        toast.success("User role updated successfully.");
        if (!users) return;
        users = users.map((user) =>
          user.ID === id ? { ...user, Role: newRole } : user,
        );
      })
      .catch((error) => sendError("Error updating user role:", error));
  }

  function deleteUser(id: string) {
    axios
      .delete(`/api/users/${id}`)
      .then(() => {
        toast.success("User deleted successfully.");
        if (!users) return;
        users = users.filter((user) => user.ID !== id);
      })
      .catch((error) => sendError("Error deleting user:", error));
  }

  function deleteSession(id: string) {
    axios
      .delete(`/api/users/${id}/session`)
      .then(() => {
        toast.success("User session deleted successfully.");
      })
      .catch((error) => sendError("Error deleting user session:", error));
  }
</script>

<h1 class="flex items-center gap-4">
  <UserCog />
  Users
</h1>

<p class="mb-4 text-muted-foreground">
  Users are automatically created when they log in with Github. You can manage
  their roles and sessions from here.
</p>

<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-[300px]">username</Table.Head>
      <Table.Head>github id</Table.Head>
      <Table.Head>role</Table.Head>
      <Table.Head>actions</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#if users === null}
      <Table.Row>
        <Table.Cell
          colspan={4}
          class="text-center animate-pulse text-muted-foreground"
        >
          Loading users...
        </Table.Cell>
      </Table.Row>
    {:else if users.length === 0}
      <Table.Row>
        <Table.Cell colspan={4} class="text-center text-muted-foreground"
          >No users found.</Table.Cell
        >
      </Table.Row>
    {:else}
      {#each users as user (user.GithubID)}
        <Table.Row>
          <Table.Cell class="flex gap-2 items-center">
            <Avatar.Root class="h-6 w-6 rounded-full">
              <Avatar.Image
                class="rounded-full"
                src={"https://avatars.githubusercontent.com/u/" +
                  user.GithubID +
                  "?v=4"}
                alt={user.Name + " avatar"}
              />
              <Avatar.Fallback class="rounded-full">{user.Name}</Avatar.Fallback
              >
            </Avatar.Root>
            {user.Name}</Table.Cell
          >
          <Table.Cell>{user.GithubID}</Table.Cell>
          <Table.Cell>{user.Role}</Table.Cell>
          <Table.Cell>
            <DropdownMenu.Root bind:open={user.Open}>
              <DropdownMenu.Trigger>...</DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Group>
                  <DropdownMenu.Label>{user.Name}</DropdownMenu.Label>
                  <DropdownMenu.Separator />
                  {#if user.Role !== "admin"}
                    <DropdownMenu.Item
                      onclick={() => changeRole(user.ID, "admin")}
                      >Make admin</DropdownMenu.Item
                    >
                  {/if}
                  {#if user.Role !== "user"}
                    <DropdownMenu.Item
                      onclick={() => changeRole(user.ID, "user")}
                      >Make user</DropdownMenu.Item
                    >
                  {/if}
                  {#if user.Role !== "guest"}
                    <DropdownMenu.Item
                      onclick={() => changeRole(user.ID, "guest")}
                      >Make guest</DropdownMenu.Item
                    >
                  {/if}
                  <DropdownMenu.Separator />
                  <DropdownMenu.Item onclick={() => deleteSession(user.ID)}
                    >Delete session</DropdownMenu.Item
                  >
                  <AlertDialog.Root>
                    <AlertDialog.Trigger class="w-full">
                      <DropdownMenu.Item closeOnSelect={false}
                        >Delete user</DropdownMenu.Item
                      >
                    </AlertDialog.Trigger>
                    <AlertDialog.Content>
                      <AlertDialog.Header>
                        <AlertDialog.Title
                          >Are you absolutely sure?</AlertDialog.Title
                        >
                        <AlertDialog.Description>
                          This action cannot be undone. This will permanently
                          delete the account and remove the data from the
                          servers.
                        </AlertDialog.Description>
                      </AlertDialog.Header>
                      <AlertDialog.Footer>
                        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
                        <AlertDialog.Action
                          onclick={() => {
                            deleteUser(user.ID);
                            user.Open = false;
                          }}>Continue</AlertDialog.Action
                        >
                      </AlertDialog.Footer>
                    </AlertDialog.Content>
                  </AlertDialog.Root>
                </DropdownMenu.Group>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </Table.Cell>
        </Table.Row>
      {/each}
    {/if}
  </Table.Body>
</Table.Root>
