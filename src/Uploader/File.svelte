<script>
  import { createEventDispatcher } from 'svelte';
  import Rusha from 'rusha';
  import Upload from './Upload.js';
  import Meta from './Meta.js';
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
  let progress = [];
  let speed = 0;
  let updates = 0;
  let percent = 0;

  let u = Upload();
  u.on('start', onUploadStart);
  u.on('finish', onUploadFinish);
  u.on('progress', onUploadProgress);
  u.on('abort', onUploadAbort);

  // any time progress changes, this runs.
  $: {
    updates = progress.length;
    let first = progress[0];
    let last = progress[progress.length-1];

    if (first != undefined && last != undefined) {
      percent = last.loaded / last.total;

      let elapsedTime = (last.time - first.time) / 1000;
      let sizeDiff = last.loaded - first.loaded;
      let bytesPerSecondStr = bytes(sizeDiff / elapsedTime);
      speed = `${bytesPerSecondStr}/s`;
    }
  }

  // actions
  function start() {
    progress = []; // clear progress history, if any.
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
    progress = [...progress, e];
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

  function hashFile() {
    status = STATUS_HASHING;
    let w = Rusha.createWorker();
    w.onmessage = (e) => {
      file.hash = e.data.hash;
      fetchMeta(file);
    };

    w.postMessage({id: 'doo', data: file.data});
  }

  // immediately request a hash
  hashFile();
</script>

<div>
  <pre>
  {JSON.stringify(file, null, 2)}
  </pre>
  status: {status}

  {#if status == STATUS_IN_PROGRESS}
    <pre>percent: {percent}</pre>
    <pre>speed: {speed}</pre>
    <pre>updates: {updates}</pre>
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
