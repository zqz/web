<script lang="ts">
	import FileListEntry from "./FileListEntry.svelte";
	import type { Uploadable } from "$lib/types";

  export let files : Array<Uploadable> = [];
</script>

<div class="file-list">
  {#each files as file (file.internalId)}
    <FileListEntry file={file} on:file:remove />
  {/each}
</div>

<style>
  .file-list{
    display: flex;
    flex-direction: column-reverse;
  }
</style>
