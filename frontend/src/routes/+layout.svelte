<script lang="ts">
import "../app.css";
import { version } from '$app/environment'
import { ModeWatcher } from 'mode-watcher';
import LightSwitch from "$lib/components/LightSwitch.svelte";
import { Toaster } from "$lib/components/ui/sonner";

</script>

<svelte:head>
</svelte:head>

<ModeWatcher />
<Toaster/>

<div class="flex flex-col min-w-[500px] h-full font-normal">
  <div class="flex-grow flex flex-col w-full p-8">
    <slot />
  </div>

  <footer class="p-8 text-sm text-center text-slate-500 font-small drop-shadow-md">
    <p>2012 - 2024</p>
    <p>{version}</p>
    <LightSwitch/>
  </footer>
</div>
