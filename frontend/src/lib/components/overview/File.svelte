<script lang="ts">
import { bytes, URLs, color } from '$lib/util';

import type { Meta } from '$lib/types';
import Thumbnail from './Thumbnail.svelte';
import FileSize from './FileSize.svelte';
import FileHash from './FileHash.svelte';
import ColorIcon from './ColorIcon.svelte';

export let file : Meta;

let thumbVisible = false;
let thumbPosX : string;
let thumbPosY : string;
let entry : HTMLElement;
let renderTop = false;

function onMouseMove(e: MouseEvent) {
  const x = e.pageX;
  const y = e.pageY;

  thumbPosX = `${x}px`;
  if (e.clientY < 300) {
    thumbPosY = `${y + 20}px`;
    renderTop = false;
  } else {
    thumbPosY = `${y - 310}px`;
    renderTop = true;
  }
}

function showPreview(e: MouseEvent) {
  onMouseMove(e);
  thumbVisible = true;
  entry.addEventListener('mousemove', onMouseMove);
}

function hidePreview() {
  thumbVisible = false;
  entry.removeEventListener('mousemove', onMouseMove);
}
</script>

<div 
  bind:this={entry}
  class="flex flex-row justify-between"
  role="tooltip"
  on:mouseover={showPreview}
  on:mouseout={hidePreview}
  on:focus={()=>{}}
  on:blur={hidePreview}
>
  <div class="flex flex-row items-center font-light
    hover:font-semibold hover:tracking-tight">
    <ColorIcon file={file}/>
    <a
      class="ml-2 flex-grow"
      href={URLs.getFileBySlugUrl(file.slug)}
      target="_blank">
      {file.name}
    </a>
  </div>
  <div class="w-32 text-right">
    <FileSize size={file.size}/>
  </div>
</div>
