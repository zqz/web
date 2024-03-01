<script lang="ts">
import { URLs } from '$lib/util';
import Uploader from '$lib/components/uploader/Uploader.svelte';
import FileList from './FileList.svelte';
import Button from '$lib/components/Button.svelte';
import { onMount } from 'svelte';
import type { Meta } from '$lib/types';
import Preview from './Preview.svelte';

let page = 0;
let delay = 1;
let selectedFileId : number | null = null;

function navigateUp() {
  console.log('up');
  if (selectedFileId === null) {
    return;
  }
  if (selectedFileId > 0) {
    selectedFileId--;
  }
}

function navigateDown() {
  console.log('down');
  if (selectedFileId === null) {
    selectedFileId = 0;
    return;
  }

  if (selectedFileId === files.length - 1) {
    selectedFileId = files.length - 1;
    return
  }

  selectedFileId++;
}

function escape() {
  selectedFileId = null;
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'ArrowUp') {
    navigateUp();
  }

  if (e.key === 'ArrowDown') {
    navigateDown();
  }

  if (e.key === 'Escape') {
    escape();
  }
}

onMount(function() {
  document.addEventListener('keydown', onKeyDown);
  return () => document.removeEventListener('keydown', onKeyDown)
});

function loadNext() {
  page++;
  loadFiles();
}

let files : Array<Meta> = [];
let error : Error | undefined;

async function fetchFiles(page: number) {
  const res = await fetch(URLs.getFilesListUrl(page));
  const json = await res.json();

  if (res.ok) {
    return json;
  } else {
    throw new Error(json);
  }
}

function loadFiles() {
  console.log('fetching files');
  fetchFiles(page).then((newFiles) => {
    files = [...files, ...newFiles];
    error = undefined;
  }).catch(e => {
    error = e;
    timeoutLoadFiles();
  });
}

function timeoutLoadFiles() {
  delay = delay + 1;
  setTimeout(loadFiles, delay*1000);
}

function selectFileId(e: CustomEvent) {
  selectedFileId = e.detail;
}

onMount(function() {
  loadFiles();
});

</script>

<Uploader on:file:uploaded={loadFiles}/>

<div class="h-full flex gap-8 flex-row">
  <div class="basis-1/2">
    <FileList files={files} on:select={selectFileId} selectedFileId={selectedFileId} />
    {#if error}
      <p>
        There was an error, retrying in {delay}s... {error.message}
      </p>
    {/if}
    <Button title="load more" on:click={loadNext}>load more</Button>
  </div>
  <div class="basis-1/2">
    <Preview files={files} selectedFileId={selectedFileId} />
  </div>
</div>
