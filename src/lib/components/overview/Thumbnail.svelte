<script lang="ts">
import type { Meta } from '$lib/types';
import { URLs } from '$lib/util';

export let file : Meta;
export let visible : boolean;
export let posX: string;
export let posY: string;
export let top: string;
let imgPath: string ;

$: imgPath = URLs.thumbnailUrl(file.slug);
$: bgPos = !top ? 'top left' : 'bottom left';
</script>

{#if visible}
  <div
    class="preview"
    style={
    `left: ${posX}; top: ${posY}; background-position: ${bgPos}; background-image: url(${imgPath});`}>
  </div>
{/if}

<style lang="scss">
.preview {
  position: absolute;
  display: flex;
  flex-direction: columns;
  align-items: flex-end;
  user-select: none;
  left: 0;
  top: 0;
  width: 300px;
  height: 300px;
  background-repeat: no-repeat;
}
</style>
