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
    max-width: 100%;
  }

  @media (min-width: 768px) and (max-width: 1024px) {
    /* tablets */
    #left {
      min-width: 100px;
    }
  }

  @media (max-width: 767px) {
    /* low res, mobile landscape */
    #left {
      display: none;
      width: 0;
      min-width: 0;
    }

    #right {
      padding-left: 0;
    }
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>
