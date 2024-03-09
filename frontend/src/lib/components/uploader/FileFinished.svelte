<script lang="ts">
import type { Uploadable } from "$lib/types";
import { Button } from "$lib/components/ui/button";
import FileHash from "../overview/FileHash.svelte";
import FileSize from "../overview/FileSize.svelte";
import LinkToFileButton from "../LinkToFileButton.svelte";
import { createEventDispatcher } from 'svelte';
import FileContainer from "./FileContainer.svelte";

const dispatch = createEventDispatcher();

export let file : Uploadable;
const meta = file.meta!;

function remove() {
  dispatch('remove', file);
}
</script>

<FileContainer file={file}>
  <div slot="buttons">
    <LinkToFileButton file={meta}>view file</LinkToFileButton>
    <Button title="Remove uploaded file from list" on:click={remove}>x</Button>
  </div>
  <div>
    <FileHash hash={meta.hash}/>
  </div>
  <div class="flex flex-row">
    <FileSize size={meta.size}/>
  </div>
</FileContainer>

