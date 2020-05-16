<script>
  import { createEventDispatcher } from 'svelte';
  import Config from '../Config.js';
  import Upload from './Upload.js';
  import Meta from './Meta.js';
  import Hashing from './Hashing.svelte';
  import Progress from './Progress.svelte';
  import ProgressStats from './ProgressStats.svelte';
  import bytes from '../Util/FileSize.js';

  export let file;

  const STATUS_QUEUE = 'queued';
  const STATUS_HASHING = 'hashing';
  const STATUS_META_CHECK = 'meta_check';
  const STATUS_READY = 'ready';
  const STATUS_IN_PROGRESS = 'in_progres';
  const STATUS_DONE = 'done';

  const dispatch = createEventDispatcher();

  let status = STATUS_QUEUE;
  let progressUpdates = [];
  let percent;

  $: {
    let first = progressUpdates[0];
    let last = progressUpdates[progressUpdates.length-1];

    if (first != undefined && last != undefined) {
      percent = last.loaded / last.total;
    }
  }

  let u = Upload();
  u.on('start', onUploadStart);
  u.on('finish', onUploadFinish);
  u.on('progress', onUploadProgress);
  u.on('abort', onUploadAbort);

  // actions
  function start() {
    progressUpdates = []; // clear progress history, if any.
    u.setFile(file);

    if (file.meta !== undefined) {
      u.upload();
      return;
    }

    let m = Meta(file);
    m.on('create', function() {
      file.meta = m.get();
      u.upload();
    });

    m.create();
  }

  function cancel() {
    u.abort();
  }

  function remove() {

  }

  // meta callbacks
  function onMetaFound(m) {
    file.meta = m;

    if (m.bytes_received == m.size) {
      status = STATUS_DONE;
    } else {
      status = STATUS_READY;
    }
  }

  function onMetaNotFound() {
    status = STATUS_READY;
  }

  // upload callbacks
  function onUploadStart() {
    status = STATUS_IN_PROGRESS;
  }

  function onUploadFinish(meta) {
    file.meta = meta;
    dispatch('file:uploaded');
    status = STATUS_DONE;
  }

  function onUploadProgress(e) {
    progressUpdates = [...progressUpdates, e];
  }

  function onUploadAbort() {
    status = STATUS_QUEUE;
    fetchMeta(file);
  };

  function fetchMeta(file) {
    status = STATUS_META_CHECK;

    let m = Meta(file);
    m.on('found', onMetaFound);
    m.on('notfound', onMetaNotFound);
    m.retrieve();
  }

  async function hashFile() {
    status = STATUS_HASHING;

    var bp = file.data.arrayBuffer();
    bp.then(async (b) => {
      const hashBuffer = await crypto.subtle.digest('SHA-1', b);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
      file.hash = hashHex;
      fetchMeta(file);
    });

  }

  function copyAttributes() {
    file.name = file.data.name;
    file.size = file.data.size;
  }

  function truncate(v, len) {
    if (v.length < len - 3) {
      return v;
    }

    return `${v.substr(0, len)}...`;
  }

  // immediately request a hash
  copyAttributes();
  hashFile();
</script>

<div>
  <div class="file {status}">
  <div class="row">
    <div class="name">
      {truncate(file.name, 20)}
    </div>
    <div class="buttons">
      {#if status === STATUS_DONE}
        <a href={`${Config.url}/api/file/by-slug/${file.meta.slug}`}>Go to
          File</a>
      {/if}
      buttons
    </div>
  </div>
  <div class="row">
    <div class="stats">
      <div>
        {bytes(file.size)}
      </div>
      <div class="monospace">
      {#if status === STATUS_HASHING}
        <Hashing/>
      {/if}
      {#if file.hash !== undefined}
        {file.hash}
      {/if}
      </div>
    </div>
    <div class="status">
      {#if status === STATUS_IN_PROGRESS}
        <ProgressStats updates={progressUpdates}/> 
      {:else}
        {status}
      {/if}
    </div>
  </div>
  {#if status === STATUS_IN_PROGRESS}
    <Progress percent={percent} height={6}/>
  {/if}
</div>
<div>
  <pre>
  {JSON.stringify(file, null, 2)}
  </pre>
  status: {status}

  {#if status === STATUS_HASHING}
    <Hashing/>
  {/if}

  {#if file.hash !== undefined && status == STATUS_READY}
  <button on:click={start}>
    start
  </button>
  <button on:click={remove}>
    remove
  </button>
  {/if}
  {#if status == STATUS_IN_PROGRESS}
    <button on:click={cancel}>
      cancel
    </button>
  {/if}
</div>
</div>

<style type="text/scss">
  @import "../variables.scss";
  .file {
    display: flex;
    flex-direction: column;

    border-left: solid 4px $link-normal;
    padding-left: 4px;

    &.done {
      border-left: solid 4px $file-done;
    }
    &:hover {
      border-width: 8px;
    }

    .row {
      font-size: 0.8rem;

      .name {
        font-size: 1.3rem;
        font-weight: 500;
        padding-bottom: 4px;
      }
    }
  }
</style>
