<script>
  import UploadFile from './UploadFile.svelte';
  export let files;
</script>

{#each files as file}
  <UploadFile file={file} on:file:remove/>
{/each}
