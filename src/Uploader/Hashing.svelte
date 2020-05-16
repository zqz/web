<script>
  let ticks = 0;
  let text = '';

  function getText(ticks) {
    const repititions = ticks % 4;
    const dots = ".".repeat(repititions);
    return "hashing" + dots;
  }

  $: text = getText(ticks);
  function update() {
    ticks++;
    setTimeout(update, 500);
  }
  update();
</script>

{text}
