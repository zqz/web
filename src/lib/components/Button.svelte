<script lang="ts">
export let title: string;
export let size = '';
</script>

<button aria-label={title} on:click class="bg-black text-white rounded-sm">
    <slot></slot>
</button>
