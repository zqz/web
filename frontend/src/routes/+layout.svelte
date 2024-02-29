<script lang="ts">
import "../app.css";
import Logo from "./Logo.svelte";
import ThemeToggle from "./ThemeToggle.svelte";
import { version } from '$app/environment'
</script>

<svelte:head>
</svelte:head>

<div class="md:px-4 full font-normal">
  <div class="flex w-full flex-col gap-8 p-8 md:flex-row md:p-16">
    <div class="md:shrink-0">
      <Logo/>
      <ThemeToggle/>
    </div>

    <div class="">
      <slot />
    </div>
  </div>

  <footer>
    <div class="text-center font-small drop-shadow-md">
      2012 - 2024<br/>
      {version}
    </div>
  </footer>
</div>
