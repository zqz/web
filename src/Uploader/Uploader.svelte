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
    document.addEventListener('drop', onDrop);
    document.addEventListener('dragover', onDragOver);
    document.addEventListener('selectFiles', openFileSelect);
  });

  let b = new Blob(['edaed'], {type: 'text/plain'});
  b.name = 'awdawd awdpokwadawd awd awawdawd awd awddexample';
  let files = [{
    id: 1,
    data: b,
  }];

  let container;

  function openFileSelect() {
    uploader.click();
  }

  function onDragOver(e) {
    e.preventDefault();
    console.log('dragover', e);
  }

  function onPaste(e) {
    console.log('paste', e);
    handleFiles(e.clipboardData.items);
  }

  function onDrop(e) {
    e.preventDefault();
    console.log('drop', e);

    debugger;
    handleFiles(e.dataTransfer.items);
  }

  function handleFiles(files) {
    let newFiles = Array.from(files)
      .map((i) => i.getAsFile())
      .filter(x => x)
      .map((i) => ({id: 1, data: i}));

    files = [...files, ...newFiles];
  }

  function onChange(e) {
    const filesToAdd = Array.from(e.target.files).map((f) => ({id: 1, data: f}));
    files = [...files, ...filesToAdd];
  }
</script>

<div class="file-list" bind:this={container}>
  {#each files as file}
    <File file={file} on:file:uploaded />
  {/each}
</div>

<style>
  .file-list{
    display: flex;
    flex-direction: column-reverse;
  }
</style>
