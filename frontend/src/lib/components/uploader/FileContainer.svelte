<script lang="ts">
import { truncate } from "$lib/text";
import type { Uploadable } from "$lib/types";
import Divider from "./Divider.svelte";

export let file: Uploadable | undefined;

const fullName = (file && file.data) ? file.data.name : 'loading';
const name = truncate(fullName, 40);
</script>

{#if file === undefined}
<div>loading</div>
{:else}
<div class="container flex flex-col rounded-md border border-solid
  border-slate-200 bg-white p-3 shadow-xl shadow-slate-300/50">
  <div class="flex flex-row justify-between">
    <div class="basis-3/4 text-xl font-bold">
      {name}
    </div>
    <div class="basis-1/4 text-right">
      <slot name="buttons"/>
    </div>
  </div>
  <Divider />
  <slot/>
</div>
{/if}
