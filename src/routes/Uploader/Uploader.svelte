<script lang="ts">
	import { onMount } from 'svelte';
  import File from './File.svelte';
  import type { Uploadable } from './types';
  import { generateId } from '$lib/text';

  let uploader: HTMLInputElement;

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

  let files : Array<Uploadable> = [];
  let container : HTMLElement;

  function onDragOver(e: DragEvent) {
    e.preventDefault();
  }

  function onPaste(e: ClipboardEvent) {
    // https://microsoft.github.io/PowerBI-JavaScript/interfaces/_node_modules_typedoc_node_modules_typescript_lib_lib_dom_d_.datatransfer.html
    if (e.clipboardData === undefined || e.clipboardData === null) {
      return;
    }

    handleFiles(e.clipboardData.items);
  }

  function onDrop(e: DragEvent) {
    e.preventDefault();
    if (e.dataTransfer === undefined || e.dataTransfer === null) {
      return;
    }
    
    handleFiles(e.dataTransfer.items);
  }

  function handleFiles(x: DataTransferItemList) {
    let newFiles = Array.from(x)
      .map((i) => i.getAsFile())
      .filter(x => x)
      .map((i) => ({id: generateId(), data: i!}));

    files = [...files, ...newFiles];
  }

  function openFileSelect() {
    uploader.click();
  }

  function onChange(e: Event) {
    const target = e.target as HTMLInputElement;

    if (target.files === undefined || target.files === null) {
      return;
    }

    if (target.files.length === 0) {
      return;
    }
    
    const filesToAdd = Array.from(target.files).map((f) => ({id: generateId(), data: f}));
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
