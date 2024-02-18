<script lang="ts">
import Overview from '$lib/components/overview/Index.svelte';
import AddFilesButton from './AddFilesButton.svelte';
import Page from './Page.svelte';
</script>

<svelte:head>
  <title>zqz.ca</title>
  <meta name="description" content="" />
</svelte:head>

<Page>
  <div slot="buttons">
    <AddFilesButton />
  </div>

  <Overview/>
</Page>
