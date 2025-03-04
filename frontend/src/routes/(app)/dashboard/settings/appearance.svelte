<script lang="ts">
  import { setMode, resetMode } from "mode-watcher";
  import Label from "$lib/components/ui/label/label.svelte";
  import Button from "$lib/components/ui/button/button.svelte";
</script>

<h2 class="text-xl font-semibold">Theme</h2>

<div class="grid grid-cols-2 gap-5 max-w-3xl">
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
    <span class="block w-full p-2 text-center font-normal"> Light </span>
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
    <span class="block w-full p-2 text-center font-normal"> Dark </span>
  </Label>
</div>
<Button class="text-muted-foreground" variant="link" onclick={() => resetMode()}
  >or system theme</Button
>
