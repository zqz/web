<script lang="ts">
  import Button from '$lib/components/Button.svelte';

  function onClick() {
    var e = new Event('selectFiles');
    document.dispatchEvent(e);
  }
</script>

<Button title="Add file to upload queue" className="" on:click={onClick}>
  add files
</Button>
