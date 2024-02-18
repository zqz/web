<script lang="ts">
import type { Meta, Uploadable } from "$lib/types";
import { URLs, truncate } from "$lib/util";
import Button from "./Button.svelte";
import Close from "../Close.svelte";
import LinkButton from "../LinkButton.svelte";
import FileHash from "../overview/FileHash.svelte";
import FileSize from "../overview/FileSize.svelte";
import LinkToFileButton from "../LinkToFileButton.svelte";
import { createEventDispatcher } from 'svelte';

const dispatch = createEventDispatcher();

export let file : Uploadable;
const meta = file.meta!;
const name = truncate(meta.name ?? 'noname' , 40);
function remove() {
  dispatch('remove', file);
}
</script>

<div class="flex flex-col rounded-md border border-slate-500 drop-shadow-md
  bg-white my-3 p-2 pt-2">
  <div class="flex flex-row justify-between">
    <div class="flex flex-col">
      <div class="text-md uppercase tracking-tight">
        {name}
      </div>
      <FileHash hash={meta.hash}/>
    </div>


    <div class="flex flex-row group">
      <LinkToFileButton file={meta}>
        view file111
      </LinkToFileButton>
      <Button title="Remove uploaded file from list" on:click={remove}>
        x
      </Button>
    </div>
  </div>
  <div class="flex flex-row">
    <FileSize size={meta.size}/>
  </div>
</div>

