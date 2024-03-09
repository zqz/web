<script lang="ts">
import type { FileProgress, Uploadable } from "$lib/types";
import { Button } from "$lib/components/ui/button";
import FileHash from "../overview/FileHash.svelte";
import FileSize from "../overview/FileSize.svelte";
import { createEventDispatcher } from 'svelte';
import FileContainer from "./FileContainer.svelte";
import ProgressStats from "./ProgressStats.svelte";
import Progress from './Progress.svelte';
import { calcPercent } from './percent';

const dispatch = createEventDispatcher();

export let file : Uploadable;
export let updates: Array<FileProgress> = [];
const meta = file.meta!;
let percent: number = 0;

$: {
  percent = calcPercent(updates);
}

function cancel() {
  dispatch('cancel', file);
}
</script>

<FileContainer file={file}>
  <div class="flex justify-end gap-1" slot="buttons">
    <Button title="Cancel file upload" size="sm" on:click={cancel}>
      Cancel
    </Button>
  </div>
  <div class="flex flex-row">
    <div class="basis-3/4">
      <FileHash hash={meta.hash}/>
    </div>
    <div class="basis-1/4 text-right text-black">
      <FileSize size={file.data.size}/>
    </div>
  </div>
  <div>
    <ProgressStats updates={updates}/>
    <Progress percent={percent}/>
  </div>
</FileContainer>
