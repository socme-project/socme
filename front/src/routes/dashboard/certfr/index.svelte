<script lang="ts">
  import axios from "axios";
  import { onMount } from "svelte";
  import * as Tabs from "$lib/components/ui/tabs/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { ShieldCheck } from "@lucide/svelte";
  import { sendError } from "$src/lib/utils";

  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";

  let alerts: any = $state([]);
  let cti: any = $state([]);
  let avis: any = $state([]);

  onMount(async () => {
    await axios
      .get("/api/misc/certfr")
      .then((res) => {
        alerts = res.data.alerts;
        cti = res.data.cti;
        avis = res.data.avis;
      })
      .catch((err) => {
        sendError("Failed to fetch Cert-FR data", err);
      });
  });
</script>

<h1 class="flex items-center gap-4 mb-2">
  <ShieldCheck />
  Cert-FR
</h1>
<p class="text-muted-foreground mb-4">
  This page displays the latest alerts, CTI, and Avis from Cert-FR.
</p>

<Tabs.Root value="alerts" class="w-full">
  <Tabs.List>
    <Tabs.Trigger value="alerts">Alerts</Tabs.Trigger>
    <Tabs.Trigger value="cti">CTI</Tabs.Trigger>
    <Tabs.Trigger value="avis">Avis</Tabs.Trigger>
  </Tabs.List>

  <Tabs.Content value="alerts">
    <h2>Alerts</h2>
    <div class="flex flex-col gap-5">
      {#if alerts.length === 0}
        {#each Array(5) as _}
          <Card.Root>
            <Card.Header>
              <Card.Title class="font-medium">
                <Skeleton class="w-3/5 h-8 mb-1" />
                <Skeleton class="w-2/3 h-4" />
              </Card.Title>
            </Card.Header>
            <Card.Content class="flex flex-col gap-1 opacity-70">
              <Skeleton class="w-4/5 h-4" />
              <Skeleton class="w-5/5 h-4" />
              <Skeleton class="w-3/5 h-4" />
            </Card.Content>
          </Card.Root>
        {/each}
      {:else}
        {#each alerts as alert}
          <Card.Root>
            <Card.Header>
              <a href={alert.Link} target="_blank">
                <Card.Title class="font-medium">{alert.Title}</Card.Title>
              </a>
              <Card.Description class="text-muted"
                >{alert.Ref} | {alert.Date} | {alert.Status}</Card.Description
              >
            </Card.Header>
            <Card.Content class="text-muted-foreground">
              {alert.Description}
            </Card.Content>
          </Card.Root>
        {/each}
      {/if}
    </div>
  </Tabs.Content>

  <Tabs.Content value="cti">
    <h2>CTI</h2>
    <div class="flex flex-col gap-5">
      {#if cti.length === 0}
        {#each Array(5) as _}
          <Skeleton />
        {/each}
      {:else}
        {#each cti as c}
          <Card.Root>
            <Card.Header>
              <a href={c.Link} target="_blank">
                <Card.Title class="font-medium">{c.Title}</Card.Title>
              </a>
              <Card.Description class="text-muted"
                >{c.Ref} | {c.Date}</Card.Description
              >
            </Card.Header>
            <Card.Content class="text-muted-foreground">
              {c.Description}
            </Card.Content>
          </Card.Root>
        {/each}
      {/if}
    </div>
  </Tabs.Content>

  <Tabs.Content value="avis">
    <h2>Avis</h2>
    <div class="flex flex-col gap-5">
      {#if avis.length === 0}
        {#each Array(5) as _}
          <Skeleton />
        {/each}
      {:else}
        {#each avis as avi}
          <Card.Root>
            <Card.Header>
              <a href={avi.Link} target="_blank">
                <Card.Title class="font-medium">{avi.Title}</Card.Title>
              </a>
              <Card.Description class="text-muted"
                >{avi.Ref} | {avi.Date}</Card.Description
              >
            </Card.Header>
            <Card.Content class="text-muted-foreground">
              {avi.Description}
            </Card.Content>
          </Card.Root>
        {/each}
      {/if}
    </div>
  </Tabs.Content>
</Tabs.Root>
