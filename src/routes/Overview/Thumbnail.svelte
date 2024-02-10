<script>
  import URLs from '$lib/urls';
  export let file;
  export let visible;
  export let posX;
  export let posY;
  export let top;
  let imgPath;
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
    img {
      user-select: none;
      border-radius: 3px;
      max-width: 300px;
      max-height: 300px;
    }
  }
</style>
