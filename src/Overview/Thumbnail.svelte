<script>
  import Config from '../Config.js';
  export let file;
  export let visible;
  export let posX;
  export let posY;
  let imgPath;
  $: imgPath = Config.thumbnailUrl(file.slug);
</script>

{#if visible}
  <div
    class="preview"
    style={`left: ${posX}; top: ${posY}`}>
    <img src={imgPath} />
  </div>
{/if}

<style type="text/scss">
  .preview {
    position: absolute;
    user-select: none;
    left: 0;
    top: 0;
    img {
      user-select: none;
      border-radius: 3px;
      max-width: 300px;
      max-height: 300px;
    }
  }
</style>
