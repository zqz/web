<script lang="ts">
  import { bytes, URLs, color } from '$lib/util';

  import type { Meta } from '$lib/types';
  import Thumbnail from './Thumbnail.svelte';

  export let file : Meta;
  const mediaType = extractMediaType(file);

  function extractMediaType(file: Meta) {
    const parts = file.type.split('/');
    return parts[0];
  }

  function isImage() : boolean {
    return mediaType === 'image';
  }

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
    if (!isImage()) {
      return;
    }

    onMouseMove(e);
    thumbVisible = true;
    entry.addEventListener('mousemove', onMouseMove);
  }

  function hidePreview() {
    if (!isImage()) {
      return;
    }
    thumbVisible = false;
    entry.removeEventListener('mousemove', onMouseMove);
  }
</script>

<div 
  bind:this={entry}
  class="entry"
  role="tooltip"
  on:mouseover={showPreview}
  on:mouseout={hidePreview}
  on:focus={()=>{}}
  on:blur={hidePreview}
  >
  <div class="row">
    <Thumbnail top={renderTop} visible={thumbVisible} posX={thumbPosX} posY={thumbPosY} file={file}/>
    <div class="sq" style="background-color: {color(mediaType)}"></div>
    <a
      class="name"
      href={URLs.getFileBySlugUrl(file.slug)}
      target="_blank">
      {file.name}
    </a>
  </div>
  <span class="monospace hash">{file.hash}</span>
  <span class="size small">{bytes(file.size)}</span>
</div>

<style lang="scss">
  @import "$lib/variables.scss";

  :global(body.dark-mode) {
    .entry {
      color: $white;

      a {
        color: $white;
      }

      &:hover {
        background-color: darken(#434a54, 10);
      }
    }
  }

  .entry {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-bottom: 4px;
    color: $black;

    .row {
      align-items: center;
    }

    &:hover {
      background-color: #e6e9ed;
    }

    .sq {
      min-width: 8px;
      width: 8px;
      height: 100%;
    }

    a {
      padding-left: 8px;
      text-decoration: none;

      &:hover {
        font-weight: 500;
      }
    }

    .hash {
      display: none;
      width: 10%;
      flex-grow: 1;
      text-align: right;
      align-self: center;
      font-size: 0.85rem;
    }

    .small {
      align-self: center;
      font-size: 0.8rem;
    }

    .size {
      width: 16%;
      text-align: right;
    }
  }

  @media (min-width: 1200px) {
    /* tablets */
    .entry {
      .hash {
        display: block;
      }
    }
  }

  @media (max-width: 767px) {
    .entry {
      .row {
        overflow:hidden;
        max-width: 70%;
        .name {
          white-space: nowrap;
          overflow:hidden;
        }
      }
    }
  }
</style>
