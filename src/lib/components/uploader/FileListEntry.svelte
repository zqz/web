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
import FileSize from '../overview/FileSize.svelte';
import FileFinished from './FileFinished.svelte';
import FileContainer from './FileContainer.svelte';
import { calcPercent } from './percent';
import Divider from './Divider.svelte';

export let file: Uploadable;

const dispatch = createEventDispatcher();

let status = FileStatus.Queue;
let progressUpdates: Array<FileProgress> = [];
let percent: number = 0;

$: {
  percent = calcPercent(progressUpdates);
}

const fileUrl = file.meta ? URLs.getFileBySlugUrl(file.meta.slug) : "invalid";

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

{#if status == FileStatus.Done}
  <FileFinished file={file} on:remove={remove}/>
{:else}
  <FileContainer file={file}>
    <div slot="buttons">
      {#if status == FileStatus.Ready}
        <Button title="Start uploading" on:click={start}>
          {#if file.meta && file.meta.bytes_received > 0}
            continue
          {:else}
            start
          {/if}
        </Button>
        <Button title="Remove file from queue" on:click={remove}>
          x
        </Button>
      {/if}
      {#if status === FileStatus.InProgress}
        <Button title="Cancel file upload" on:click={cancel}>
          cancel
        </Button>
      {/if}
    </div>
    
    <div class="flex flex-row">
      <div class="basis-3/4">
        {#if status === FileStatus.Hashing}
          <Hashing/>
        {/if}
        {#if file.hash !== undefined}
          {file.hash}
        {/if}
      </div>
      <div class="basis-1/4 text-right">
        <FileSize size={file.data.size}/>
        <span class="monospace">{status}</span>
      </div>
    </div>

    {#if status === FileStatus.InProgress}
      <div>
        <ProgressStats updates={progressUpdates}/>
        <Progress percent={percent}/>
      </div>
    {/if}
  </FileContainer>
{/if}
