<script>
  import bytes from './Util/FileSize.js';
  export let file;

  const url = `https://x.zqz.ca/${file.slug}`;
</script>

<div class="entry">
  <a href={url} target="_blank">{file.name}</a>
  <span>{bytes(file.size)}</span>
</div>

<style type="text/scss">
  .entry {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-bottom: 4px;

    a {
      text-decoration: none;
      color: #4A89DC;

      &:visited {
        color: #967ADC;
      }

      &:hover {
        font-weight: 500;
      }
    }

  }
</style>
