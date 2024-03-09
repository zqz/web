<script lang="ts">
import type { Meta } from '$lib/types';
import File from './File.svelte';
	import { createEventDispatcher } from 'svelte';

export let files : Array<Meta>;
export let selectedFileId : number | null;

const dispatch = createEventDispatcher();

function onClick(id: number) {
  console.log('on click', id);
  dispatch('select', id);
}
</script>

<div class="flex flex-col mb-3 overflow-hidden">
  {#each files as file, i (file.slug)}
    <a on:click={() => onClick(i)}>
    <File file={file} selected={i == selectedFileId}/>
    </a>
  {/each}
</div>
