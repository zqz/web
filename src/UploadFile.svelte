<script>
  import { createEventDispatcher } from 'svelte';
import Rusha from 'rusha';
import Upload from './Upload.js';
  export let file;

  const dispatch = createEventDispatcher();
  const STATUS_QUEUE = 'queued';
  const STATUS_IN_PROGRESS = 'in_progres';

  let status = STATUS_QUEUE;

  function start() {
    status = STATUS_IN_PROGRESS;
    
    let u = Upload(file);
    u.onStart = (e) => {
      console.log(e);
    }

    u.onChunkStart = (e) => {}
    u.onChunkFinish = (e) => {}
    u.onFinish = (e) => {}
    u.onProgress = (e) => {}
    u.upload();
  }

  function hashFile() {
    if (hasHash()) {
      return;
    }

    let w = Rusha.createWorker();
    w.onmessage = (e) => {
      file.hash = e.data.hash;
    };

    w.postMessage({id: 'doo', data: 'xxx'});
  }

  function cancel() {
    status = STATUS_QUEUE;
  }

  function remove() {
    dispatch('file:remove', {
      id: file.id,
    });
  }

  function hasHash() {
    return file.hash !== undefined;
  }

  hashFile();
</script>

<div>
  FILE: {JSON.stringify(file)} status: {status}
  {#if file.hash !== undefined && status == STATUS_QUEUE}
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
