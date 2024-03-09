<script lang="ts">
import type { Meta } from "$lib/types";
import { URLs } from "$lib/util";
import LinkToFileButton from "../LinkToFileButton.svelte";
import FileHash from "./FileHash.svelte";
import FileSize from "./FileSize.svelte";
import * as Card from "$lib/components/ui/card";
export let file : Meta | null;

let imgPath = "";

$: hasFile = file !== null;
$: {
  if (hasFile) {
    imgPath = URLs.thumbnailUrl(file!.slug);
  }
}

</script>


{#if file !== null}
  <Card.Root class="w-full lg:basis-1/2">
    <Card.Header>
      <div class="flex flex-row justify-between">
        <div class="text-xl font-bold">
          {file.name}
        </div>
      </div>
    </Card.Header>

    <Card.Content>
      <div>
        <img src={imgPath} alt={file.name} class="rounded-md shadow-md
          max-h-[600px]"/>
        <div class="pt-4 pb-4">
          Size: <FileSize size={file.size}/><br/>
          SHA: <FileHash hash={file.hash}/>
          Date: {file.date}<br>
        </div>
      </div>
      <LinkToFileButton file={file}>Download</LinkToFileButton>
    </Card.Content>
  </Card.Root>
{/if}
