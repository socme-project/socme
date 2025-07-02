<script lang="ts">
  import { onMount } from "svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import { user } from "$lib/stores/user";
  import FullscreenLoading from "$src/lib/components/fullscreen-loading.svelte";
  import { navigate } from "sv-router/generated";

  onMount(async () => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get("code");
    const state = params.get("state");

    if (!code) {
      console.error("OAuth code not found.");
      toast.error("Authentication failed. Please try again.");
      navigate("/");
      return;
    }

    await axios
      .get(`/api/auth/github/callback?code=${code}&state=${state}`)
      .then((res) => {
        user.set(res.data.user);
        toast.success("Authenticated successfully.");
        navigate("/dashboard");
      })
      .catch(() => {
        toast.error("Failed to authenticate.");
        user.set(null);
        navigate("/");
      });
  });
</script>

<FullscreenLoading />
