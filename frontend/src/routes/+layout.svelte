<script lang="ts">
import "../app.css";
import { version } from '$app/environment'
</script>

<svelte:head>
</svelte:head>

<div class="flex flex-col md:px-4 h-full font-normal">
  <div class="flex-grow flex flex-col w-full p-8">
    <slot />
  </div>

  <footer class="p-8 text-sm text-center text-slate-500 font-small drop-shadow-md">
    2012 - 2024 # {version}
  </footer>
</div>
