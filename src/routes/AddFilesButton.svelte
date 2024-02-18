<script lang="ts">
  import Button from '$lib/components/uploader/Button.svelte';

  function onClick() {
    var e = new Event('selectFiles');
    document.dispatchEvent(e);
  }
</script>

<Button title="Add file to upload queue" className="text-xl p-3" on:click={onClick}>
  add files
</Button>
