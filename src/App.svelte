<script>
  import Overview from './Overview/Index.svelte';
  import Logo from './Logo.svelte';
  import AddFilesButton from './AddFilesButton.svelte';
</script>

<main>
  <div id="left">
    <Logo/>
  </div>
  <div id="right">
    <header>
      <h1>zqz.ca</h1>
      <AddFilesButton />
    </header>
    <Overview/>
  </div>
</main>

<style type="text/scss">
  @import "variables.scss";

  main {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }


  header {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
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
    color: $black;
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
