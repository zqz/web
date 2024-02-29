<script lang="ts">
  export let percent: number;

  let width = 0;
  $: {
    width = percent * 100;
  }
</script>

<div class="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700">
  <div class="bg-blue-600 h-2.5 rounded-full" style="width: {width}%"></div>
</div>
