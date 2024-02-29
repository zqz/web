<script lang="ts">
import type { Meta } from '$lib/types';
import { color } from '$lib/util';
export let file : Meta;
const mediaType = extractMediaType(file);
  function isImage() : boolean {
    return mediaType === 'image';
  }

function extractMediaType(file: Meta) {
  const parts = file.type.split('/');
  return parts[0];
}

const iconColor = color(mediaType);
</script>
<div class="h-5 min-w-2 bg-blue-400 inline-block" style="background-color: {iconColor}"></div>
