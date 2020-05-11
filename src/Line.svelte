<script>
  export let height;
  export let color;

  const lineStyle = function() {
    return `height: ${height}px; width: ${height}px; background-color: ${color};`;
  }
</script>

<div class="line-container">
  <div class="line-start" style={lineStyle()}></div>
  <div class="line-fill" style={lineStyle()}></div>
  <div class="line-end" style={lineStyle()}></div>
</div>

<style type="text/scss">
  .line-container {
    display: flex;
    flex-direction: row;
    margin-bottom: 8px;
  }

  .line-start {
    border-top-left-radius: 50%;
    border-bottom-left-radius: 50%;
  }
  .line-end {
    border-top-right-radius: 50%;
    border-bottom-right-radius: 50%;
  }
  .line-fill {
    flex-grow: 1;
  }
</style>
