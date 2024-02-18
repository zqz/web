
<script lang="ts">
	import type { Meta } from '$lib/types';
  import File from './File.svelte';

  export let files : Array<Meta>;
</script>

<div class="flex flex-col mb-3">
  {#each files as file (file.hash)}
    <File file={file}/>
  {/each}
</div>
