<script lang="ts">
import type { Uploadable } from "$lib/types";
import { Button } from "$lib/components/ui/button";
import FileSize from "../overview/FileSize.svelte";
import FileContainer from "./FileContainer.svelte";
import Close from '../Close.svelte';
import { createEventDispatcher } from 'svelte';

const dispatch = createEventDispatcher();

export let file : Uploadable;

const isStarted = file.meta && file.meta.bytes_received > 0;
const bytesSeen = file.meta?.bytes_received ?? 0;
const buttonText = isStarted ? 'Resume' : 'Start';
const buttonTitle = isStarted ? 'Resume Uploading' : 'Start Upload';

function remove() {
  dispatch('remove', file);
}

function start() {
  dispatch('start', file);
}
</script>

<FileContainer file={file}>
  <div class="flex justify-end gap-1" slot="buttons">
    <Button title={buttonTitle} size="sm" on:click={start}>
      {buttonText}
    </Button>
    <Button title="Remove file from queue" size="sm" on:click={remove}>
      <Close/>
    </Button>
  </div>
  <div class="flex flex-row">
    <div class="basis-3/4">
      {file.hash}
    </div>
    <div class="basis-1/4 text-right text-black">
      {#if isStarted}
        <FileSize size={bytesSeen}/> of
      {/if}
      <FileSize size={file.data.size}/>
    </div>
  </div>
</FileContainer>
