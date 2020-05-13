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
  });

  let files = [];
  let container;

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
