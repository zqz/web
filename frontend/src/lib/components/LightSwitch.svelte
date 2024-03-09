<script lang="ts">
  import Sun from "svelte-radix/Sun.svelte";
  import Moon from "svelte-radix/Moon.svelte";
  import { setMode, mode } from 'mode-watcher';

	function toggleMode() {
		if ($mode === 'light') {
			setMode('dark');
		} else {
			setMode('light');
		}
	}
</script>
 
<button on:click={toggleMode}> 
  <Sun
    class="h-[1rem] w-[1rem] rotate-0 scale-100 transition-all dark:hidden
    dark:-rotate-90 dark:scale-0"
  />
  <Moon
    class="h-[1rem] w-[1rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"
  />
  <span class="sr-only">Toggle theme</span>
</button>
