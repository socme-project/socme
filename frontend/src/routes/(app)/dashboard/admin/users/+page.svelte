<script lang="ts">
  import { Crown, UserCog } from "lucide-svelte";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import Ellipsis from "lucide-svelte/icons/ellipsis";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import { onMount } from "svelte";
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";

  let users: any = $state([]);

  onMount(async () => {
    try {
      const res = await axios.get("/api/users/list", {
        headers: { Authorization: localStorage.getItem("token") },
      });
      users = res.data.users;
    } catch (error) {
      toast.error("Internal server error");
    }
  });
</script>

<h1 class="flex items-center gap-4">
  <UserCog />
  Users
</h1>

<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-[300px]">username</Table.Head>
      <Table.Head>github ID</Table.Head>
      <Table.Head>Role</Table.Head>
      <Table.Head></Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#if users.length === 0}
      {#each Array(5) as _}
        <Table.Row>
          <Table.Cell class="flex justify-start items-center gap-4">
            <Skeleton class="size-12 rounded-full" />
            <Skeleton class="w-[100px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[120px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[70px] h-4" />
          </Table.Cell>
          <Table.Cell>
            <Skeleton class="w-[20px] h-3" />
          </Table.Cell>
        </Table.Row>
      {/each}
    {:else}
      {#each users as user}
        <Table.Row>
          <Table.Cell class="flex justify-start items-center gap-4">
            <Avatar.Root>
              <Avatar.Image src={user.Avatar} alt={user.Username + "logo"} />
              <Avatar.Fallback>
                <Ellipsis />
              </Avatar.Fallback>
            </Avatar.Root>
            <div class="flex items-center gap-2">
              {#if user.Role === "admin"}
                <Crown class="text-yellow-700 dark:text-yellow-100" size={16} />
              {/if}
              {user.Username}
            </div></Table.Cell
          >
          <Table.Cell>{user.GitHubID}</Table.Cell>
          <Table.Cell>{user.Role}</Table.Cell>

          <Table.Cell class="text-right w-10">
            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                <Ellipsis size={18} />
              </DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Group>
                  <DropdownMenu.GroupHeading
                    >{user.Username}</DropdownMenu.GroupHeading
                  >
                  <DropdownMenu.Separator />
                  {#if user.Role !== "guest"}
                    <DropdownMenu.Item
                      onclick={() => {
                        let role = "guest";
                        axios
                          .get("/api/user/change-role", {
                            headers: {
                              Authorization: localStorage.getItem("token"),
                            },
                            params: { id: user.ID, role: role },
                          })
                          .then((res) => {
                            if (res.status === 200) {
                              toast.success("Role changed to " + role);
                            } else {
                              toast.error("Error:" + res.data.message);
                            }
                          })
                          .catch((error) => {
                            toast.error("Error:" + error);
                          });
                      }}>Make guest</DropdownMenu.Item
                    >
                  {/if}
                  {#if user.Role !== "user"}
                    <DropdownMenu.Item
                      onclick={() => {
                        let role = "user";
                        axios
                          .get("/api/user/change-role", {
                            headers: {
                              Authorization: localStorage.getItem("token"),
                            },
                            params: { id: user.ID, role: role },
                          })
                          .then((res) => {
                            if (res.status === 200) {
                              toast.success("Role changed to " + role);
                            } else {
                              toast.error("Error:" + res.data.message);
                            }
                          })
                          .catch((error) => {
                            toast.error("Error:" + error);
                          });
                      }}>Make user</DropdownMenu.Item
                    >
                  {/if}
                  {#if user.Role !== "admin"}
                    <DropdownMenu.Item
                      onclick={() => {
                        let role = "admin";
                        axios
                          .get("/api/user/change-role", {
                            headers: {
                              Authorization: localStorage.getItem("token"),
                            },
                            params: { id: user.ID, role: role },
                          })
                          .then((res) => {
                            if (res.status === 200) {
                              toast.success("Role changed to " + role);
                            } else {
                              toast.error("Error:" + res.data.message);
                            }
                          })
                          .catch((error) => {
                            toast.error("Error:" + error);
                          });
                      }}>Make admin</DropdownMenu.Item
                    >
                  {/if}

                  <DropdownMenu.Item
                    onclick={() => {
                      axios
                        .get("/api/user/revoke-session", {
                          headers: {
                            Authorization: localStorage.getItem("token"),
                          },
                          params: { id: user.ID },
                        })
                        .then((res) => {
                          if (res.status === 200) {
                            toast.success("Session revoked");
                          } else {
                            toast.error("Error:" + res.data.message);
                          }
                        })
                        .catch((error) => {
                          toast.error("Error:" + error);
                        });
                    }}>Revoke Session</DropdownMenu.Item
                  >

                  <DropdownMenu.Separator />
                  <DropdownMenu.Item
                    onclick={() => {
                      axios
                        .get("/api/user/delete", {
                          headers: {
                            Authorization: localStorage.getItem("token"),
                          },
                          params: { id: user.ID },
                        })
                        .then((res) => {
                          if (res.status === 200) {
                            toast.success("User deleted");
                          } else {
                            toast.error("Error:" + res.data.message);
                          }
                        })
                        .catch((error) => {
                          toast.error("Error:" + error);
                        });
                    }}>Delete User</DropdownMenu.Item
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
