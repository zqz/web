<script lang="ts">
import { onMount } from 'svelte';

const maxTicks = 4;
const updateFrequency = 500;
let ticks = 0;
let text = '';

$: text = getText(ticks);

function getText(ticks: number) {
  const repititions = ticks % maxTicks;
  const dots = ".".repeat(repititions);
  return "hashing" + dots;
}

onMount(() => {
  const interval = setInterval(() => {
    ticks++;
    console.log('still ticking');
    }, updateFrequency);

  return () => clearInterval(interval);
});
</script>

<div>{text}</div>
