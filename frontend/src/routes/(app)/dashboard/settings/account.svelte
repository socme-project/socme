<script lang="ts">
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import { PencilOff } from "lucide-svelte/icons";
  import Input from "$lib/components/ui/input/input.svelte";
  import { Label } from "$lib/components/ui/label/index.js";

  import { user } from "$lib/stores/user";
</script>

<h2>Account Informations</h2>
{#if $user}
  <div class="flex flex-wrap py-4 gap-12">
    <div>
      <div class="relative">
        <Avatar.Root
          class="h-20 w-20 md:h-40 md:w-40 rounded-full p-1 border-muted-foreground border"
        >
          <Avatar.Image
            class="rounded-full"
            src={$user.Avatar}
            alt={$user.Username + " avatar"}
          />
          <Avatar.Fallback class="rounded-full"
            >{$user.Username}</Avatar.Fallback
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
        <Input type="text" id="username" value={$user.Username} disabled />
      </div>
      <div class="flex w-full max-w-sm flex-col gap-1.5">
        <Label for="githubid">Github ID</Label>
        <Input type="text" id="githubid" value={$user.GitHubID} disabled />
      </div>
      <div class="flex w-full max-w-sm flex-col gap-1.5">
        <Label for="role">Role</Label>
        <Input type="text" id="role" value={$user.Role} disabled />
      </div>
    </div>
  </div>
  <p class="text-muted-foreground text-xs pt-3">
    (!) Account informations are in read-only mode. Change informations from the
    Auth provider (github).
  </p>
{/if}
