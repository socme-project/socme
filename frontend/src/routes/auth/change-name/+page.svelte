<script lang="ts">
  import { user } from "$lib/stores/user";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import axios from "axios";
  import { toast } from "svelte-sonner";

  let newUsername = "";
  let currentUser = $user;

  async function changeUsername() {
    if (!currentUser) {
      toast.error("You must be logged in to change your username.");
      return;
    }
    if (!newUsername.trim()) {
      toast.error("Username cannot be empty.");
      return;
    }

    try {
      const response = await axios.post(
        "/api/auth/change-name",
        { Username: newUsername },
        { withCredentials: true },
      );

      if (response.data.success) {
        user.set({
          ID: currentUser.ID,
          Username: newUsername,
          GitHubID: currentUser.GitHubID ?? "",
          GoogleID: currentUser.GoogleID ?? "",
          Role: currentUser.Role ?? "guest",
          Avatar: currentUser.Avatar ?? "",
        });
        toast.success("Username changed successfully.");
        goto("/dashboard");
      } else {
        toast.error(response.data.message || "Failed to change username.");
      }
    } catch (error) {
      console.error("Error changing username:", error);
      toast.error("An error occurred. Please try again.");
    }
  }
</script>

<div class="h-screen w-screen flex justify-center items-center">
  <div class="text-center bg-white p-6 rounded-lg shadow-lg">
    <h1 class="text-2xl font-bold">Change Your Username</h1>
    <p class="text-gray-600">Please choose a new username.</p>
    <input
      type="text"
      bind:value={newUsername}
      class="border p-2 rounded mt-4 w-full"
      placeholder="Enter new username"
    />
    <button
      on:click={changeUsername}
      class="bg-primary text-white p-2 rounded mt-4 w-full"
    >
      Save
    </button>
  </div>
</div>
