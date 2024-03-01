<script lang="ts">
import type { Meta } from '$lib/types';
import FileSize from './FileSize.svelte';
import ColorIcon from './ColorIcon.svelte';

export let file : Meta;
export let selected : boolean = false;
</script>

<div 
  class="group flex flex-row justify-between"
  role="tooltip"
  data-selected={selected}
>
  <div class="flex flex-row items-center font-light
    hover:font-semibold hover:tracking-tight
    group-data-[selected=true]:underline group-data-[selected=true]:drop-shadow">
    <ColorIcon file={file}/>
    <div class="ml-2 flex-grow">{file.name}</div>
  </div>
  <div class="w-32 text-right">
    <FileSize size={file.size}/>
  </div>
</div>
