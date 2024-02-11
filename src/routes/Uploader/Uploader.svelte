<script lang="ts">
import FileList from './FileList.svelte';
import { onMount } from 'svelte';
import type { Uploadable } from './types';
import { generateId } from '$lib/text';

let uploader: HTMLInputElement;

onMount(() => {
  document.addEventListener('paste', onPaste);
  document.addEventListener('drop', onDrop);
  document.addEventListener('dragover', onDragOver);
  document.addEventListener('selectFiles', openFileSelect);
});

let files : Array<Uploadable> = [];

function onDragOver(e: DragEvent) {
  e.preventDefault();
}

function onPaste(e: ClipboardEvent) {
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
  .map(i => i.getAsFile())
  .filter(x => x)
  .map(f => uploadableFromFile(f!));

  files = [...files, ...newFiles];
}

function openFileSelect() {
  if (uploader === undefined) {
    throw "could not find uploader on page";
  }

  uploader.click()
}

function onChange(e: Event) {
  const target = e.target as HTMLInputElement;

  if (target.files === undefined || target.files === null) {
    return;
  }

  if (target.files.length === 0) {
    return;
  }

  const filesToAdd = Array.from(target.files).map(f => uploadableFromFile(f));
  files = [...files, ...filesToAdd];
}

function onFileRemove(e: CustomEvent) {
  const message = e.detail as Uploadable;
  files = files.filter(f => f.internalId !== message.internalId);
}

function uploadableFromFile(f: globalThis.File) {
  return {
    data: f,
    internalId: generateId()
  }
}
</script>

<div>
  <input 
    bind:this={uploader}
    on:change={onChange}
    type="file" multiple
    class="hidden-upload"/>

  <FileList files={files} on:file:remove={onFileRemove} />
</div>

<style lang="scss">
.hidden-upload {
  display: none;
}
</style>
