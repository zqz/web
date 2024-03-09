<script lang="ts">
import type { Meta } from "$lib/types";
import { URLs } from "$lib/util";
import LinkToFileButton from "../LinkToFileButton.svelte";
import FileHash from "./FileHash.svelte";
import FileSize from "./FileSize.svelte";
export let file : Meta | null;

let imgPath = "";

$: hasFile = file !== null;
$: {
  if (hasFile) {
    imgPath = URLs.thumbnailUrl(file!.slug);
  }
}

</script>

<div class={"transition duration-500 ease-in-out " + (hasFile ? "opacity-100" : "opacity-0")}>
  <div class="opacity-100 p-8 rounded-md bg-white h-full drop-shadow-lg duration-200 ">
    <div class="h-full justify-between flex flex-col">
      {#if file !== null}
        <div>
          <img src={imgPath} alt={file.name} class="rounded-md shadow-md"/>
          <div class="text-xl font-light pt-4 pb-4">{file.name}</div>
          <div>
            Size: <FileSize size={file.size}/><br/>
            SHA: <FileHash hash={file.hash}/>
            Date: {file.date}<br>
          </div>
        </div>
        <LinkToFileButton file={file}>Download</LinkToFileButton>
      {/if}
    </div>
  </div>
</div>
