<script lang="ts">
</script>

<style lang="scss">
  .logo-container {
    display: flex;
    flex-direction: column;

    .logo {
      transition: transform 8s ease-in 5.5s;
    }

    &:hover {
      .logo {
        transform: rotate(8640deg);
      }
    }
  }

  .logo {
    user-select: none;
    background-image: url("https://avatars1.githubusercontent.com/u/111311?s=60&amp;v=4");
    background-size: 100%;
    width: 64px;
    height: 64px;
    outline: 0;
    border: 0;
  }
</style>

<div class="logo-container">
  <div class="logo"/>
</div>
