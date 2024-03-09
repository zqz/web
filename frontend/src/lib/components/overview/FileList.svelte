<script lang="ts">
import type { Meta } from '$lib/types';
import File from './File.svelte';
import { createEventDispatcher } from 'svelte';

import * as Table from "$lib/components/ui/table";
	import { Size } from 'svelte-radix';
	import FileSize from './FileSize.svelte';


export let files : Array<Meta>;
export let selectedFileId : number | null;

const dispatch = createEventDispatcher();

function onClick(id: number) {
  console.log('on click', id);
  dispatch('select', id);
}
</script>

<div class="flex flex-col mb-3 max-h-[calc(70vh)] overflow-scroll">
  <Table.Root class="max-h-full">
    <Table.Header>
      <Table.Row>
        <Table.Head>File</Table.Head>
        <Table.Head>Size</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body class="max-h-full">
    {#each files as file, i (file.slug)}
      <Table.Row data-state={selectedFileId === i ? "selected" : ""}
        on:click={() => onClick(i)}>
        <Table.Cell>
          {file.name}
        </Table.Cell>
        <Table.Cell>
          <FileSize size={file.size}/>
        </Table.Cell>
      </Table.Row>
      {/each}
    </Table.Body>
  </Table.Root>
</div>
