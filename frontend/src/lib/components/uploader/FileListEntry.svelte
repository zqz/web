<script lang="ts">
import { createEventDispatcher, onMount } from 'svelte';
import { uploadFile } from './uploadFile';

import Hashing from './Hashing.svelte';
import Progress from './Progress.svelte';
import ProgressStats from './ProgressStats.svelte';
import { FileEvent, FileStatus, type FileProgress, type Uploadable } from '$lib/types';
import FileSize from '../overview/FileSize.svelte';
import FileFinished from './FileFinished.svelte';
import { Button } from '$lib/components/ui/button';
import FileContainer from './FileContainer.svelte';
import { calcPercent } from './percent';

export let file: Uploadable;

const dispatch = createEventDispatcher();

let status = FileStatus.Queue;
let updates: Array<FileProgress> = [];
let percent: number = 0;

$: {
  percent = calcPercent(updates);
}

let u = uploadFile(file);
u.on(FileEvent.Error, () => status = FileStatus.Error);
u.on(FileEvent.Start, () => status = FileStatus.InProgress);
u.on(FileEvent.Progress, (e: FileProgress) => updates = [...updates, e]);
u.on(FileEvent.MetaNotFound, () => status = FileStatus.Ready);
u.on(FileEvent.Hash, () => status = FileStatus.Hashing);
u.on(FileEvent.Finish, () => {
  status = FileStatus.Done;
});

u.on(FileEvent.Abort, () => {
  status = FileStatus.Queue;
  u.fetchMeta();
});

u.on(FileEvent.MetaFound, () => {
  const m = file.meta!;

  if (m.bytes_received === m.size) {
    status = FileStatus.Done;
  } else {
    status = FileStatus.Ready;
  }
});

function start() {
  updates = [];
  u.start();
}

function cancel() {
  u.abort();
}

function remove() {
  dispatch('file:remove', file);
}

// request that the file is hashed as soon as the component is mounted.
onMount(u.hash);
</script>

{#if status == FileStatus.Done}
  <FileFinished file={file} on:remove={remove}/>
{:else}
  <FileContainer file={file}>
    <div slot="buttons">
      {#if status == FileStatus.Ready}
        <Button title="Start uploading" size="sm" on:click={start}>
          {#if file.meta && file.meta.bytes_received > 0}
            Continue
          {:else}
            Start
          {/if}
        </Button>
        <Button title="Remove file from queue" size="sm" on:click={remove}>
          x
        </Button>
      {/if}
      {#if status === FileStatus.InProgress}
        <Button title="Cancel file upload" size="sm" on:click={cancel}>
          Cancel
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
        <ProgressStats updates={updates}/>
        <Progress percent={percent}/>
      </div>
    {/if}
  </FileContainer>
{/if}
