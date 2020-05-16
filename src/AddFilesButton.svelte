<script>
  function onClick() {
    var event = new Event('selectFiles');
    document.dispatchEvent(event);
  }
</script>

<span class="add-files-button" on:click={onClick}>
  add files
</span>

<style type="text/scss">
  @import "variables.scss";

  .add-files-button {
    cursor: pointer;
    font-weight: 500;
    color: white;
    background-color: $link-normal;
    padding: 10px 20px;
    border-radius: 3px;
    flex-shrink: 0;
  }
</style>
