<script lang="ts">
	import type { FileProgress } from '$lib/types';
  import { bytes } from '$lib/util';
  export let updates : Array<FileProgress>;

  let transferred: string;
  let speed: string;
  let total: string;
  let estimate: string;

  // any time progress changes, this runs.
  $: {
    let first = updates[0];
    let last = updates[updates.length-1];

    if (first != undefined && last != undefined) {
      let elapsedTime = (last.time - first.time) / 1000;
      let sizeDiff = last.loaded - first.loaded;
      let bytesPerSecond = sizeDiff / elapsedTime;
      let timeLeftInSeconds = (last.total - last.loaded) / bytesPerSecond;
      let bytesPerSecondStr = bytes(bytesPerSecond);
      speed = `${bytesPerSecondStr}/s`;
      if (!isNaN(timeLeftInSeconds)) {
        estimate = `${timeLeftInSeconds.toFixed(0)}s left`
      }
      total = bytes(first.total);
      transferred = bytes(last.loaded);
    }
  }
</script>

<div class="flex">
  <div class="basis-1/2">
    {transferred} / {total}
  </div>
  <div class="basis-1/2 text-right">
    {estimate} @ {speed}
  </div>
</div>
