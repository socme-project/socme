<script lang="ts">
	import type { Snippet } from 'svelte';
	let { children }: { children: Snippet } = $props();
  import { ModeWatcher } from "mode-watcher";
	import { Toaster } from '$lib/components/ui/sonner/index.js';
</script>

<ModeWatcher />
<Toaster />
{@render children()}
