<script lang="ts">
  import { Button } from '$lib/components/ui/button';

  function onClick() {
    var e = new Event('selectFiles');
    document.dispatchEvent(e);
  }
</script>

<Button title="Add file to upload queue" size="sm" on:click={onClick}>
  Add files
</Button>
