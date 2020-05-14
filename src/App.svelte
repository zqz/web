<script>
  import Overview from './Overview/Index.svelte';
  import Logo from './Logo.svelte';
</script>

<main>
  <div id="left">
    <Logo/>
  </div>
  <div id="right">
    <h1>zqz.ca</h1>
    <Overview/>
  </div>
</main>

<style>
  main {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }

  h1 {
    font-weight: 500;
    font-size: 32px;
    user-select: none;
    padding: 0;
    margin: 0;
    margin-bottom: 8px;
    line-height: 48px;
    vertical-align: middle;
    width: 100%;
  }

  #left {
    min-width: 256px;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
  }

  #right {
    padding-left: 32px;
    flex-grow: 1;
    display: flex;
    flex-direction: column;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>
