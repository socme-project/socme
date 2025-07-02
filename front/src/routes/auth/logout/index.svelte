<script lang="ts">
  import { user } from "$lib/stores/user";
  import FullscreenLoading from "$src/lib/components/fullscreen-loading.svelte";
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import { SESSION_COOKIE_NAME } from "$lib/variables";
  import { navigate } from "sv-router/generated";

  function deleteCookie(cookieName: string) {
    document.cookie = `${cookieName}=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
  }

  onMount(async () => {
    // FIXME: Should we let the user know about the error? Do we log out ?
    // await axios
    //   .get(`/api/auth/logout`)
    //   .then(() => {})
    //   .catch((error) => {
    //     console.log("Error while deleting session: " + error);
    //   });

    user.set(null);
    deleteCookie(SESSION_COOKIE_NAME);
    toast.success("Logged out successfully");
    navigate("/")
  });
</script>

<FullscreenLoading />
