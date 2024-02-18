<script lang="ts">
  export let title: string;
  export let className: string;
</script>

<button aria-label={title} title={title} class="btn {className}" on:click>
  <slot></slot>
</button>
