<script lang="ts">
  export let height;
  export let percent;

  let bgColor = '#434a54';
  let fgColor = '#a0d468';
  let line;

  const endStyle = function(color) {
    return `height: ${height}px; width: ${height}px; background-color: ${color};`;
  }

  const lineStyle = function(color) {
    return `height: ${height}px; width: ${percent * 100}%; background-color: ${color};`;
  }

  const adjustStyle = function() {
    return `top: -${height}px;`;
  }

  $: {
    if (line !== undefined) {
      line.style.width = `${percent * 100}%`;
    }
  }
</script>

<div class="line-container">
  <div class="back">
    <div class="line-start" style={endStyle(fgColor)}></div>
    <div class="line-fill" style={lineStyle(bgColor)}></div>
    <div class="line-end" style={endStyle(bgColor)}></div>
  </div>
  <div class="front" style="{adjustStyle()}">
    <div class="line-start" style={endStyle(fgColor)}></div>
    <div class="line-progress" bind:this={line} style={lineStyle(fgColor)}></div>
    <div class="line-end" style={endStyle(fgColor)}></div>
  </div>
</div>

<style lang="scss">
  .back {
    display: flex;
    flex-direction: row;
  }

  .front {
    display: flex;
    justify-content: flex-start;
    position: relative;
  }

  .line-start {
    border-top-left-radius: 50%;
    border-bottom-left-radius: 50%;
  }

  .line-end {
    border-top-right-radius: 50%;
    border-bottom-right-radius: 50%;
  }

  .line-progress {
    transition: width 1s linear 0s;
  }

  .line-fill {
    flex-grow: 1;
  }
</style>
