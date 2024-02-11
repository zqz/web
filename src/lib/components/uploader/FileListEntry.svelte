<script lang="ts">
import { createEventDispatcher } from 'svelte';
import { bytes, URLs, truncate } from '$lib/util';
import { uploadFile } from './uploadFile';

import Hashing from './Hashing.svelte';
import Button from '$lib/components/Button.svelte';
import LinkButton from '$lib/components/LinkButton.svelte';
import Progress from './Progress.svelte';
import ProgressStats from './ProgressStats.svelte';
import { FileEvent, FileStatus, type FileProgress, type Uploadable } from '$lib/types';

export let file: Uploadable;

const dispatch = createEventDispatcher();

let status = FileStatus.Queue;
let progressUpdates: Array<FileProgress> = [];
let percent: number = 0;

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
  status = FileStatus.Hashing;
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
  dispatch('file:remove', file);
}

// meta callbacks
function onMetaFound() {
  console.log('meta found');
  const m = file.meta!;

  if (m.bytes_received === m.size) {
    status = FileStatus.Done;
  } else {
    status = FileStatus.Ready;
  }
}

function onMetaNotFound() {
  console.log('meta not found');
  status = FileStatus.Ready;
}

// upload callbacks
function onUploadStart() {
  status = FileStatus.InProgress;
}

function onUploadFinish() {
  console.log('upload finished');
  dispatch('file:uploaded');
  status = FileStatus.Done;
}

function onUploadProgress(e: FileProgress) {
  progressUpdates = [...progressUpdates, e];
}

function onUploadAbort() {
  status = FileStatus.Queue;
  u.fetchMeta();
};

function onMetaCheck() {
  status = FileStatus.MetaCheck;
}

// immediately request a hash
u.hash();
</script>

<div class="file {status}">
  <div class="row">
    <div class="name">
      {truncate(file.data.name, 40)}
    </div>
    <div class="row">
      {#if status == FileStatus.Ready}
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
      {#if status === FileStatus.InProgress}
        <Button title="Cancel file upload" on:click={cancel}>
          cancel
        </Button>
      {/if}
      {#if status === FileStatus.Done}
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
        {bytes(file.data.size)}
      </div>
      <div class="monospace">
        {#if status === FileStatus.Hashing}
          <Hashing/>
        {/if}
        {#if file.hash !== undefined}
          {file.hash}
        {/if}
      </div>
    </div>
    <div class="row">
      <div class="status">
        {#if status === FileStatus.InProgress}
          <ProgressStats updates={progressUpdates}/>
        {:else}
          <span class="monospace">{status}</span>
        {/if}
      </div>
    </div>
  </div>
  {#if status === FileStatus.InProgress}
    <Progress percent={percent} height={6}/>
  {/if}
</div>

<style lang="scss">
@import "$lib/variables.scss";
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
