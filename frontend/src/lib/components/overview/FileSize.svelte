<script lang="ts">
import { bytes } from "$lib/util";
export let size : number;
$: humanSize = bytes(size);
</script>

<span class="text-nowrap text-xs font-thin subpixel-antialiased">
  {humanSize} 
</span>
