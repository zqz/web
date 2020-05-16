<script>
  import Config from '../Config.js';
  import Uploader from '../Uploader/Uploader.svelte';
  import Entry from './Entry.svelte';

  let promise = fetchFiles();
  let delay = 1;

  async function fetchFiles() {
    const res = await fetch(`${Config.url}/api/files`);
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
</script>

<Uploader on:file:uploaded={loadFiles}/>
<div>
  {#await promise}
    <p>
      Loading Files...
    </p>
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
</div>
