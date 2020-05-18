<script>
  import bytes from '../Util/FileSize.js';
  import Config from '../Config.js';

  export let file;
  const [registry, name] = file.type.split('/');

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

    return '#434a54';
  }
</script>

<div class="entry">
  <div class="row">
    <div class="sq" style="background-color: {color()}"></div>
    <a class="name" href={Config.getFileBySlugUrl(file.slug)} target="_blank">{file.name}</a>
  </div>
  <span class="monospace hash">{file.hash}</span>
  <span class="size small">{bytes(file.size)}</span>
</div>

<style type="text/scss">
  @import "../variables.scss";

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
