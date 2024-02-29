<script lang="ts">
import FileListEntry from "./FileListEntry.svelte";
import type { Uploadable } from "$lib/types";

export let files : Array<Uploadable> = [];
</script>

<div class="flex flex-col-reverse gap-2 mb-3">
  {#each files as file (file.internalId)}
    <FileListEntry file={file} on:file:remove on:file:uploaded />
  {/each}
</div>
