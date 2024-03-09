<script lang="ts">
import { createEventDispatcher, onMount } from 'svelte';
import { uploadFile } from './uploadFile';

import Hashing from './Hashing.svelte';
import { FileEvent, FileStatus, type FileProgress, type Uploadable } from '$lib/types';
import FileSize from '../overview/FileSize.svelte';
import FileFinished from './FileFinished.svelte';
import FileContainer from './FileContainer.svelte';
import FileInProgress from './FileInProgress.svelte';
import FileReady from './FileReady.svelte';

export let file: Uploadable;

const dispatch = createEventDispatcher();

let status = FileStatus.Queue;
let updates: Array<FileProgress> = [];

let u = uploadFile(file);
u.on(FileEvent.Error, () => status = FileStatus.Error);
u.on(FileEvent.Start, () => status = FileStatus.InProgress);
u.on(FileEvent.Progress, (e: FileProgress) => updates = [...updates, e]);
u.on(FileEvent.MetaNotFound, () => status = FileStatus.Ready);
u.on(FileEvent.Hash, () => status = FileStatus.Hashing);
u.on(FileEvent.Finish, () => {
  status = FileStatus.Done;
  dispatch('file:uploaded', file);
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
{:else if status == FileStatus.InProgress}
  <FileInProgress file={file} updates={updates} on:cancel={cancel}/>
{:else if status == FileStatus.Ready}
  <FileReady file={file} on:start={start} on:remove={remove}/>
{:else}
  <FileContainer file={file}>
    <div class="flex justify-end gap-1" slot="buttons">
    </div>

    <div class="flex flex-row">
      <div class="basis-3/4">
        {#if status === FileStatus.Hashing}
          <Hashing/>
        {/if}
      </div>
      <div class="basis-1/4 text-right">
        <FileSize size={file.data.size}/>
        <span class="monospace">{status}</span>
      </div>
    </div>
  </FileContainer>
{/if}
