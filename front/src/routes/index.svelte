<script lang="ts">
  import Button from "$src/lib/components/ui/button/button.svelte";
  import { Github } from "@lucide/svelte";
</script>

<main>
<div class="w-full h-screen flex justify-center items-center">
  <div class="p-8">
    <div class="mx-auto flex flex-col justify-center space-y-6 w-96">
      <div class="flex flex-col space-y-2 text-center">
        <img
          src="/logo-white.png"
          alt="logo"
          width="55"
          class="hidden dark:block m-auto"
        />
        <img
          src="/logo-black.png"
          alt="logo"
          width="55"
          class="block dark:hidden m-auto"
        />
        <h1 class="text-3xl md:text-4xl font-semibold tracking-tight mb-0">
          Sign-in
        </h1>
      </div>
      <Button href={`/api/auth/github/login`} class="w-full" onclick={() => {
        window.location.href = `/api/auth/github/login`;
      }}>
        <Github /> Continue with GitHub
      </Button>
      <p class="text-muted-foreground px-8 text-center text-sm">
        By clicking continue, you agree to our
        <a href="/terms"> Terms of Service </a>
        and
        <a href="/privacy"> Privacy Policy </a>
        .
      </p>
    </div>
  </div>
</div>
</main>
