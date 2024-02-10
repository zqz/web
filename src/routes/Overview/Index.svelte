<script lang="ts">
  import Config from '../Config';
  import Uploader from '../Uploader/Uploader.svelte';
  import Entry from './Entry.svelte';
  import Button from '../Button.svelte';
  import { onMount } from "svelte";

  $: page = 0;
  $: promise = null;
  let delay = 1;

  function loadNext() {
    page++;
    promise = fetchFiles();
  }

  async function fetchFiles() {
    const res = await fetch(Config.getFilesListUrl(page));
    const json = await res.json();

    if (res.ok) {
      return json;
    } else {
      throw new Error(json);
    }
  }

  function loadFiles() {
    promise = fetchFiles();
  }

  function timeoutLoadFiles() {
    delay = delay + 1;
    setTimeout(loadFiles, delay*1000);
    return '';
  }

  onMount(() => {
    promise = fetchFiles();
  });
</script>

<div>
  <Uploader on:file:uploaded={loadFiles}/>
  <div class="files-list">
    {#if promise}
    {#await promise}
      <p>Loading Files...</p>
    {:then files}
      {#each files as f}
        <Entry file={f}/>
      {/each}

    {:catch error}
      <p>
        {timeoutLoadFiles()}
        There was an error, retrying in {delay}s...
      </p>
    {/await}
    {/if}
    <Button on:click={loadNext} size="large">load more</Button>
  </div>
</div>

<style>
  .files-list {
    margin-top: 16px;
  }
</style>
