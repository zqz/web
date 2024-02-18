<script lang="ts">
import type { Meta } from "$lib/types";
import { URLs } from "$lib/urls";
import Button from "./uploader/Button.svelte";

export let file: Meta;
const fileUrl = URLs.getFileBySlugUrl(file.slug);
</script>

<a href={fileUrl} target="_blank">
  <Button title={file.name}><slot></slot></Button>
</a>
