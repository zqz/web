<style type="text/scss">

  .logo-container {
    display: flex;
    flex-direction: column;

    .logo {
      transition: transform 8s ease-in 5.5s;
    }

    .foo {
      text-align: right;
      transition: opacity 0.5s ease-in 0s;
      opacity: 0;
    }

    &:hover {
      .logo {
        transform: rotate(8640deg);
      }
      .foo {
        opacity: 1;
      }
    }
  }


  .logo {
    background-image: url("https://avatars1.githubusercontent.com/u/111311?s=60&amp;v=4");
    background-size: 100%;
    width: 64px;
    height: 64px;

  }
</style>

<div class="logo-container">
  <div class="logo"></div>
  <div class="foo">awdad</div>
</div>
