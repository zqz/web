<script>
  import UploadList from './UploadList.svelte';

let files = [{id: '1'}];
  let x = 0;

  function handleFileAdd() {
    files = [...files, {id: x++}];
  }

  function handleFileRemoved(file) {
    files = files.filter(f => f.id != file.detail.id);
  }
</script>

<div>
  <button on:click={handleFileAdd}>
    Add
  </button>
  <div>Upload</div>

  {#if files.length > 0}
    <UploadList
      files={files}
      on:file:remove={handleFileRemoved}/>
    files {files.length}
  {:else}
    no files
  {/if}
</div>
