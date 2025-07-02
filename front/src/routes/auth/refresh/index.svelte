<script lang="ts">
  import { onMount } from "svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";
  import { user } from "$lib/stores/user";
  import FullscreenLoading from "$src/lib/components/fullscreen-loading.svelte";
  import { navigate } from "sv-router/generated";
    import { SESSION_COOKIE_NAME } from "$src/lib/variables";


  function deleteCookie(cookieName: string) {
    document.cookie = `${cookieName}=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
  }

  onMount(async () => {
    const params = new URLSearchParams(window.location.search);
    const next = params.get("next");

    await axios
      .get("/api/auth/refresh")
      .then((res) => {
        user.set(res.data.user);
        if (next) {
          navigate(next);
        } else {
          navigate("/dashboard");
        }
      })
      .catch(() => {
        toast.error("Failed to refresh session.");
        user.set(null);
        deleteCookie(SESSION_COOKIE_NAME)
        navigate("/");
      });
  });
</script>

<FullscreenLoading />
