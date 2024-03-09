<script lang="ts">
import type { Meta } from '$lib/types';
import type { SvelteComponent, SvelteComponentTyped } from 'svelte';
import Image from "svelte-radix/Image.svelte";
import SpeakerLoud from "svelte-radix/SpeakerLoud.svelte";
import Video from "svelte-radix/Video.svelte";
import File from "svelte-radix/File.svelte";
import FontStyle from "svelte-radix/FontStyle.svelte";
import FileText from "svelte-radix/FileText.svelte";

const iconClasses = "inline h-[1rem] w-[1rem]";

export let file : Meta;

function extractMediaType(file: Meta) {
  const parts = file.type.split('/');
  return parts[0];
}

function fileIcon(type: string) : typeof File {
  switch(type) {
    case 'application': return File;
    case 'audio': return SpeakerLoud;
    case 'font': return FontStyle;
    case 'image': return Image;
    case 'multipart': return File;
    case 'text': return FileText;
    case 'video': return Video;
  }
  return File;
}

const type = extractMediaType(file);
const IconComponent = fileIcon(type);

</script>

<IconComponent title={type ?? file.type} class={iconClasses}/>
