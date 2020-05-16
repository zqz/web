<script>
  import bytes from '../Util/FileSize.js';
  import Config from '../Config.js';
  export let file;

  const url = `${Config.url}/api/file/by-slug/${file.slug}`;

const [registry, name] = file.type.split('/'); // split mime
function color() {
  switch(registry) {
    case 'application': return '#434a54';
    case 'audio': return '#fcbb42';
    case 'font': return '#37bc9b';
    case 'example': return '#da4453';
    case 'image': return '#4a89dc';
    case 'message': return '#da4453';
    case 'multipart': return '#da4453';
    case 'text': return '#967adc';
    case 'video': return '#a0d468';
  }

  return '434a54';
}
</script>

<div class="entry">
  <div class="row">
    <div class="sq" style="background-color: {color()}"></div>
    <a href={url} target="_blank">{file.name}</a>
  </div>
  <span class="small">{bytes(file.size)}</span>
</div>

<style type="text/scss">
  @import "../variables.scss";

  .entry {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-bottom: 4px;

    .row {
      align-items: center;
    }

    .sq {
      width: 8px;
      height: 100%;
    }

    a {
      padding-left: 8px;
      text-decoration: none;
      color: $link-normal;

      &:hover {
        font-weight: 500;
      }
    }

    .small {
      font-size: 0.8rem;
      align-self: center;
    }
  }
</style>
