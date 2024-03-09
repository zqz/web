<script lang="ts">
import type { Uploadable } from "$lib/types";
import { Button } from "$lib/components/ui/button";
import FileHash from "../overview/FileHash.svelte";
import FileSize from "../overview/FileSize.svelte";
import LinkToFileButton from "../LinkToFileButton.svelte";
import { createEventDispatcher } from 'svelte';
import FileContainer from "./FileContainer.svelte";
import Close from "$lib/components/Close.svelte";

const dispatch = createEventDispatcher();

export let file : Uploadable;
const meta = file.meta!;

function remove() {
  dispatch('remove', file);
}
</script>

<FileContainer file={file}>
  <div class="flex justify-end gap-1" slot="buttons">
    <LinkToFileButton file={meta}>View file</LinkToFileButton>
    <Button title="Remove uploaded file from list" size="sm" on:click={remove}>
      <Close/>
    </Button>
  </div>
  <div>
    <FileHash hash={meta.hash}/>
  </div>
  <div class="flex flex-row">
    <FileSize size={meta.size}/>
  </div>
</FileContainer>

