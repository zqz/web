<script lang="ts">
  import Button from '$lib/components/Button.svelte';
  export let url = '';
  export let target = '';
  export let title = '';
</script>

<a href={url} target={target}>
  <Button title={title}>
    <slot></slot>
  </Button>
</a>

<style>
  a {
    text-decoration: none;
    outline: 0;
  }
</style>
