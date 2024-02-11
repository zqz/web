<script lang="ts">
import { createEventDispatcher } from 'svelte';
import URLs from '$lib/urls';
import { uploadFile } from './Upload';
import bytes from '$lib/size';
import { truncate } from '$lib/text';

import Hashing from './Hashing.svelte';
import Button from '../Button.svelte';
import LinkButton from '../LinkButton.svelte';
import Progress from './Progress.svelte';
import ProgressStats from './ProgressStats.svelte';
import { FileEvent, type FileProgress } from './types';

export let file : File

const STATUS_QUEUE = 'queued';
const STATUS_HASHING = 'hashing';
const STATUS_META_CHECK = 'meta_check';
const STATUS_READY = 'ready';
const STATUS_IN_PROGRESS = 'in_progres';
const STATUS_DONE = 'done';

const dispatch = createEventDispatcher();

let status = STATUS_QUEUE;
let progressUpdates: Array<FileProgress> = [];
let percent;

$: {
  let first = progressUpdates[0];
  let last = progressUpdates[progressUpdates.length-1];

  if (first != undefined && last != undefined) {
    percent = last.loaded / last.total;
  }
}

let u = uploadFile(file);
// maybe handle error here as well
u.on(FileEvent.Start, onUploadStart);
u.on(FileEvent.Finish, onUploadFinish);
u.on(FileEvent.Progress, onUploadProgress);
u.on(FileEvent.Hash, onHash);
u.on(FileEvent.Abort, onUploadAbort);
u.on(FileEvent.MetaCheck, onMetaCheck);
u.on(FileEvent.MetaFound, onMetaFound);
u.on(FileEvent.MetaNotFound, onMetaNotFound);

function onHash() {
  status = STATUS_HASHING;
}

// actions
function start() {
  progressUpdates = []; // clear progress history, if any.
  u.start();
}

function cancel() {
  u.abort();
}

function remove() {
  dispatch('file:removed', { id: file.id });
}

// meta callbacks
function onMetaFound(m) {
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

function onUploadProgress(e: FileProgress) {
  progressUpdates = [...progressUpdates, e];
}

function onUploadAbort() {
  status = STATUS_QUEUE;
  u.fetchMeta();
};

function onMetaCheck() {
  status = STATUS_META_CHECK;
}

function copyAttributes() {
  file.name = file.data.name;
  file.size = file.data.size;
}

// immediately request a hash
copyAttributes();
u.hash();
</script>

<div class="file {status}">
  <div class="row">
    <div class="name">
      {truncate(file.name, 40)}
    </div>
    <div class="row">
      {#if status == STATUS_READY}
        <Button title="Start uploading" on:click={start}>
          {#if file.meta && file.meta.bytes_received > 0}
            continue
          {:else}
            start
          {/if}
        </Button>
        <Button title="Remove file from queue" size="remove" on:click={remove}>
          x
        </Button>
      {/if}
      {#if status === STATUS_IN_PROGRESS}
        <Button title="Cancel file upload" on:click={cancel}>
          cancel
        </Button>
      {/if}
      {#if status === STATUS_DONE}
        <LinkButton title="View file" target="_blank" url={URLs.getFileBySlugUrl(file.meta.slug)}>
          goto :file
        </LinkButton>
        <Button title="Remove uploaded file from list" size="remove" on:click={remove}>
          x
        </Button>
      {/if}
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
    <div class="row">
      <div class="status">
        {#if status === STATUS_IN_PROGRESS}
          <ProgressStats updates={progressUpdates}/>
        {:else}
          <span class="monospace">{status}</span>
        {/if}
      </div>
    </div>
  </div>
  {#if status === STATUS_IN_PROGRESS}
    <Progress percent={percent} height={6}/>
  {/if}
</div>

<style lang="scss">
@import "../variables.scss";
.file {
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;

  border-left: solid 8px $link-normal;
  padding-left: 8px;
  padding-bottom: 4px;

  &.done {
    border-color: $file-done;
  }

  .row {
    font-size: 0.8rem;

    .status {
      align-self: flex-end;
    }

    .name {
      font-size: 1.3rem;
      font-weight: 500;
      padding-bottom: 4px;
    }
  }
}
</style>
