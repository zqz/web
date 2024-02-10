<script lang="ts">
	import { onMount } from 'svelte';
  import File from './File.svelte';

  let uploader = null;

  onMount(() => {
    uploader = document.createElement('input');
    uploader.type = 'file';
    uploader.addEventListener('change', onChange);
    uploader.style.display = 'none';
    uploader.multiple = true;

    container.appendChild(uploader);
    document.addEventListener('paste', onPaste);
    document.addEventListener('drop', onDrop);
    document.addEventListener('dragover', onDragOver);
    document.addEventListener('selectFiles', openFileSelect);
  });

  let files = [];
  let container;

  function onDragOver(e) {
    e.preventDefault();
  }

  function onPaste(e) {
    // https://microsoft.github.io/PowerBI-JavaScript/interfaces/_node_modules_typedoc_node_modules_typescript_lib_lib_dom_d_.datatransfer.html
    handleFiles(e.clipboardData.items);
  }

  function onDrop(e) {
    e.preventDefault();
    if (e.dataTransfer === undefined) {
      return;
    }
    handleFiles(e.dataTransfer.items);
  }

  function randomId() {
    return Math.random().toString(20).substr(2, 8)
  }

  function handleFiles(x: DataTransferItems) {
    let newFiles = Array.from(x)
      .map((i) => i.getAsFile())
      .filter(x => x)
      .map((i) => ({id: randomId(), data: i}));

    files = [...files, ...newFiles];
  }

  function openFileSelect() {
    uploader.click();
  }

  function onChange(e) {
    const filesToAdd = Array.from(e.target.files).map((f) => ({id: randomId(), data: f}));
    files = [...files, ...filesToAdd];
  }

  function onFileRemoved(e) {
    files = files.filter(x => x.id != e.detail.id);
  }
</script>

<div class="file-list" bind:this={container}>
  {#each files as file}
    <File file={file} on:file:uploaded on:file:removed={onFileRemoved} />
  {/each}
</div>

<style>
  .file-list{
    display: flex;
    flex-direction: column-reverse;
  }
</style>
