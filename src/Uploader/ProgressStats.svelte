<script>
  import bytes from '../Util/FileSize.js';
  export let updates;

  let percent;
  let speed;
  let transferred;
  let total;
  let estimate = '';

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

<div>
  <div class="row">
    {transferred} / {total}
  </div>
  <div class="row">
    {estimate} @ {speed}
  </div>
</div>

<style>
  .row {
    justify-content: flex-end;
  }
</style>
