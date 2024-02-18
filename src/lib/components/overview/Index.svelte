<script lang="ts">
  import { URLs } from '$lib/util';
  import Uploader from '$lib/components/uploader/Uploader.svelte';
  import FileList from './FileList.svelte';
  import Button from '$lib/components/Button.svelte';

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

<Uploader on:file:uploaded={loadFiles}/>

<div>
  {#await fetchFiles(page, reload)}
    <p>Loading Files...</p>
  {:then files}
    <FileList files={files}/>
  {:catch error}
    <p>
      {timeoutLoadFiles()}
      There was an error, retrying in {delay}s... {error.message}
    </p>
  {/await}
  <Button title="load more" on:click={loadNext} size="large">load more</Button>
</div>
