<script lang="ts">
  import Settings from "lucide-svelte/icons/settings";
  import * as Tabs from "$lib/components/ui/tabs/index.js";
  import Appearance from "./appearance.svelte";
  import Account from "./account.svelte";
</script>

<h1 class="flex items-center gap-4">
  <Settings />
  Settings
</h1>

<Tabs.Root value="account" class="w-full">
  <Tabs.List
    class="w-full justify-start gap-2 bg-background [&>*]:rounded-none [&>*]:border-primary"
  >
    <Tabs.Trigger class="data-[state=active]:border-b-2" value="account"
      >Account</Tabs.Trigger
    >
    <Tabs.Trigger class="data-[state=active]:border-b-2" value="appearance"
      >Appearance</Tabs.Trigger
    >
  </Tabs.List>
  <Tabs.Content value="account">
    Make changes to your account here.
    <Account />
  </Tabs.Content>
  <Tabs.Content value="appearance">
    Change the appearance of your account here.
    <Appearance />
  </Tabs.Content>
</Tabs.Root>
