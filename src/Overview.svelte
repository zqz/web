<script>
  import Entry from './Entry.svelte';

  const fetchFiles = (async() => {
    const response = await fetch('https://zqz.ca/api/files');
    return await response.json();
  })()
</script>

{#if false}
{#await fetchFiles}
in progress
{:then files}
  {#each files as f}
    <Entry file={f}/>
  {/each}
{:catch error}
error
{/await}
{/if}
