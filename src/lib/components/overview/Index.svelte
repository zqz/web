<script lang="ts">
  import { URLs } from '$lib/util';
  import Uploader from '$lib/components/uploader/Uploader.svelte';
  import Entry from './Entry.svelte';
  import Button from '$lib/components/Button.svelte';
  import { onMount } from "svelte";

  $: page = 0;
  $: reload = 0;
  let delay = 1;

  function loadNext() {
    page++;
  }

  async function fetchFiles(page: number, reload: number) {
    console.log('reload', reload);
    const res = await fetch(URLs.getFilesListUrl(page));
    const json = await res.json();

    if (res.ok) {
      return json;
    } else {
      throw new Error(json);
    }
  }

  function loadFiles() {
    console.log('reloading files');
    reload++;
    page = 0;
  }

  function timeoutLoadFiles() {
    delay = delay + 1;
    setTimeout(loadFiles, delay*1000);
    return '';
  }
</script>

<div>
  <Uploader on:file:uploaded={loadFiles}/>

  <div class="files-list">
    {#await fetchFiles(page, reload)}
      <p>Loading Files...</p>
    {:then files}
      {#each files as file (file.hash)}
        <Entry file={file}/>
      {/each}

    {:catch error}
      <p>
        {timeoutLoadFiles()}
        There was an error, retrying in {delay}s... {error.message}
      </p>
    {/await}
    <Button title="load more" on:click={loadNext} size="large">load more</Button>
  </div>
</div>

<style>
  .files-list {
    margin-top: 16px;
  }
</style>
