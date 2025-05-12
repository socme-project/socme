<script lang="ts">
  import CirclePlus from "lucide-svelte/icons/circle-plus";
  import Check from "lucide-svelte/icons/check";
  import { cn } from "$lib/utils.js";
  import * as Popover from "$lib/components/ui/popover";
  import * as Command from "$lib/components/ui/command";
  import Button from "$lib/components/ui/button/button.svelte";
  import Badge from "$lib/components/ui/badge/badge.svelte";
  import Separator from "$lib/components/ui/separator/separator.svelte";

  type Props = {
    title: string;
    selectedValue: string;
    options: {
      label: string;
      value: string;
      // This should be `Component` after lucide-svelte updates types
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      icon?: any;
    }[];
  };

  let {
    title,
    options,
    selectedValue = $bindable<string>(""),
  }: Props = $props();
</script>

<Popover.Root>
  <Popover.Trigger>
    {#snippet child({ props })}
      <Button {...props} variant="outline" size="sm" class="h-8 border-dashed">
        <CirclePlus />
        {title}
        {#if selectedValue != ""}
          <Separator orientation="vertical" class="mx-2 h-4" />
          {#each options as option}
            {#if option.value === selectedValue}
              <Badge class="rounded-sm px-1 font-normal">
                {option.label}
              </Badge>
            {/if}
          {/each}
        {/if}
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content class="w-[200px] p-0" align="start">
    <Command.Root>
      <Command.Input placeholder={title} />
      <Command.List>
        <Command.Empty>No results found.</Command.Empty>
        <Command.Group>
          {#each options as option}
            {@const isSelected = selectedValue === option.value}
            <Command.Item
              onSelect={() => {
                if (isSelected) {
                  selectedValue = "";
                } else {
                  selectedValue = option.value;
                }
              }}
            >
              <div
                class={cn(
                  "border-primary mr-2 flex size-4 items-center justify-center rounded-full border",
                  isSelected
                    ? "bg-primary text-primary-foreground"
                    : "opacity-50 [&_svg]:invisible",
                )}
              >
                <Check class="size-4" />
              </div>
              {#if option.icon}
                {@const Icon = option.icon}
                <Icon class="text-muted-foreground" />
              {/if}

              <span>{option.label}</span>
            </Command.Item>
          {/each}
        </Command.Group>
        {#if selectedValue != ""}
          <Command.Separator />
          <Command.Group>
            <Command.Item
              class="justify-center text-center"
              onSelect={() => (selectedValue = "")}
            >
              Clear filter
            </Command.Item>
          </Command.Group>
        {/if}
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
