<script lang="ts">
import { URLs } from '$lib/util';
import Uploader from '$lib/components/uploader/Uploader.svelte';
import FileList from './FileList.svelte';
import Button from '$lib/components/Button.svelte';
import { onMount } from 'svelte';
import type { Meta } from '$lib/types';

$: page = 0;
let delay = 1;

function loadNext() {
  page++;
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
    files = [...newFiles];
    error = undefined;
  }).catch(e => {
    error = e;
    console.log('error', e);
    timeoutLoadFiles();
  });
}

function timeoutLoadFiles() {
  delay = delay + 1;
  setTimeout(loadFiles, delay*1000);
  return '';
}

onMount(function() {
  loadFiles();
});

</script>

<Uploader on:file:uploaded={loadFiles}/>

<div>
  <FileList files={files}/>

  {#if error}
    <p>
      There was an error, retrying in {delay}s... {error.message}
    </p>
  {/if}

  <Button title="load more" on:click={loadNext}>load more</Button>
</div>
