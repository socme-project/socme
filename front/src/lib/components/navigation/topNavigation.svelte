<script lang="ts">
  import { Github } from "@lucide/svelte";
  import Commandbar from "./commandbar.svelte";
  import Separator from "../ui/separator/separator.svelte";
  import Button from "../ui/button/button.svelte";
  import SidebarTrigger from "../ui/sidebar/sidebar-trigger.svelte";
</script>

<header class="flex h-14 shrink-0 items-center gap-2">
  <div class="flex flex-1 items-center gap-2 px-3">
    <SidebarTrigger />
    <Separator orientation="vertical" class="mr-2 h-4" />

    <Commandbar />
    <Separator orientation="vertical" class="mr-2 h-4" />
    <Button
      type="button"
      href="https://github.com/socme-project"
      variant="ghost"
      size="icon"
      class="h-7 w-7"
    >
      <Github />
      <span class="sr-only">Github repo</span>
    </Button>
  </div>
</header>
