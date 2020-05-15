<script>
	import { onMount } from 'svelte';
  import File from './File.svelte';

  const uploader = document.createElement('input');
  uploader.type = 'file';
  uploader.addEventListener('change', onChange);
  uploader.style.display = 'none';
  uploader.multiple = true;

  onMount(() => {
    container.appendChild(uploader);
    document.addEventListener('paste', onPaste);
  });

  let files = [];
  let container;

  function onPaste(e) {
    let pastedFiles = Array.from(e.clipboardData.items)
      .map((i) => i.getAsFile())
      .filter(x => x)
      .map((i) => ({id: 1, data: i}));

    files = [...files, ...pastedFiles];
  }


  function onChange(e) {
    const filesToAdd = Array.from(e.target.files).map((f) => ({id: 1, data: f}));
    files = [...files, ...filesToAdd];
  }

  function addFile() {
    uploader.click();
  }
</script>

<div>
  <button on:click={addFile}>
    Add
  </button>

  <div bind:this={container}>
    {#each files as file}
      <File file={file} on:file:uploaded />
    {/each}
  </div>
</div>

<style>
</style>
