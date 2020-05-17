<script>
  import Button from './Button.svelte';

  function onClick() {
    var e = new Event('selectFiles');
    document.dispatchEvent(e);
  }
</script>

<Button title="Add file to upload queue" size="large" on:click={onClick}>
  add files
</Button>
