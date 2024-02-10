<script lang="ts">
  let ticks = 0;
  let text = '';

  $: text = getText(ticks);

  function getText(ticks) {
    const repititions = ticks % 4;
    const dots = ".".repeat(repititions);
    return "hashing" + dots;
  }

  function update() {
    ticks++;
    setTimeout(update, 500);
  }

  update();
</script>

{text}
