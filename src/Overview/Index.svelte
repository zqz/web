<script>
  import Config from '../Config.js';
  import Uploader from '../Uploader/Uploader.svelte';
  import Entry from './Entry.svelte';

  let promise = fetchFiles();

  async function fetchFiles() {
    const res = await fetch(`${Config.url}/api/files`);
    const json = await res.json();

    if (res.ok) {
      return json;
    } else {
      throw new Error(json);
    }
  }

  function onFileUploaded() {
    promise = fetchFiles();
  }
</script>

<Uploader on:file:uploaded={onFileUploaded}/>
{#await promise}
in progress
{:then files}
  {#each files as f}
    <Entry file={f}/>
  {/each}
{:catch error}
error
{/await}
