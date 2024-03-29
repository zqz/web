<script lang="ts">
import { URLs } from '$lib/util';
import Uploader from '$lib/components/uploader/Uploader.svelte';
import FileList from './FileList.svelte';
import { Button } from '$lib/components/ui/button';
import { onMount } from 'svelte';
import type { Meta } from '$lib/types';
import Preview from './Preview.svelte';
import * as Card from "$lib/components/ui/card";
import { toast } from 'svelte-sonner';

let page = 0;
let delay = 1;
let selectedFileId : number | null = null;
let selectedFile : Meta | null = null;

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
  loadFiles(true);
}

let files : Array<Meta> = [];
let error : Error | undefined;
const seenHashes = new Set<string>();

function addFiles(newFiles: Array<Meta>) {
  const filesToAdd : Array<Meta> = [];
  newFiles.forEach(f => {
    if (seenHashes.has(f.hash)) {
      return;
    }

    seenHashes.add(f.hash);
    filesToAdd.push(f);
  });

  files = [...filesToAdd, ...files];
}

async function fetchFiles(page: number) {
  const res = await fetch(URLs.getFilesListUrl(page));
  const json = await res.json();

  if (res.ok) {
    return json;
  } else {
    throw new Error(json);
  }
}

function loadFiles(showToast: boolean) {
  console.log('fetching files');
  fetchFiles(page).then((newFiles) => {
    addFiles(newFiles);
    error = undefined;

    if(showToast) {
      toast("Loaded more files");
    }
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
$: {
  if (selectedFileId != null) {
    selectedFile = files[selectedFileId];
  }
}

onMount(function() {
  loadFiles(false);
});

function onFileUploaded(event: CustomEvent) {
  const f = event.detail.data as Meta;
  toast(`${f.name} was uploaded successfully`)
  loadFiles(false);
}

</script>

<Uploader on:file:uploaded={onFileUploaded}/>
<div class="h-full flex gap-8 flex-column xl:flex-row">
  <Card.Root class="w-full lg:basis-1/2">

    <Card.Header>
      <div class="flex flex-row justify-between">
        <div class="basis-3/4 text-xl font-bold">
          files
        </div>
        <div class="basis-1/4 text-right">
        </div>
      </div>
    </Card.Header>

    <Card.Content>
      <FileList files={files} on:select={selectFileId} selectedFileId={selectedFileId} />
      {#if error}
        <p>
          There was an error, retrying in {delay}s... {error.message}
        </p>
      {/if}
      <Button title="load more" on:click={loadNext}>load more</Button>
    </Card.Content>
  </Card.Root>
  <div class="opacity-0 hidden lg:block lg:basis-1/2 lg:opacity-100 transition duration-500 ease-in-out">
    <Preview file={selectedFile}/>
  </div>
</div>
