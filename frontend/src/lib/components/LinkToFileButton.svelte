<script lang="ts">
import type { Meta } from "$lib/types";
import { URLs } from "$lib/urls";
import { Button } from "$lib/components/ui/button";

export let file: Meta;
$: fileUrl = URLs.getFileBySlugUrl(file.slug);
</script>

<a class="flex-inline" href={fileUrl} target="_blank">
  <Button title={file.name}>
    <slot></slot>
  </Button>
</a>
