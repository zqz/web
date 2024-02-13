<script lang="ts">
  import Settings from './Settings';

  let text = Settings.getTheme();
  let light = true;

  if (text == 'dark') {
    light = false;
    window.document.body.classList.toggle('dark-mode')
  }

  function toggle() {
    light = !light;
    text = light ? 'light' : 'dark';
    Settings.setTheme(text);
    window.document.body.classList.toggle('dark-mode')
  }
</script>

<style>
  .toggle {
    user-select: none;
    cursor: pointer;
    width: 64px;
    text-align: center;
    margin-top: 10px;
  }

  .toggle:hover {
    font-weight: bold;
  }
</style>

<button class="toggle monospace" on:click={toggle}>
  {text}
</button>
