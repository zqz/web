<script lang="ts">
import "../app.css";
	import Logo from "./Logo.svelte";
	import ThemeToggle from "./ThemeToggle.svelte";
</script>

<svelte:head>
  <link
    href="https://fonts.googleapis.com/css2?family=Inter:400,100,200,300,500,600,700,800,900&display=swap" rel="stylesheet">
  <link href="https://fonts.googleapis.com/css2?family=Fira+Code&display=swap" rel="stylesheet">
</svelte:head>

<div class="md:px-4 full font-normal">
  <div class="flex w-full flex-col gap-8 p-8 md:flex-row md:p-16">
    <div class="md:shrink-0">
      <Logo/>
      <ThemeToggle/>
    </div>

    <div class="gap-8 flex-grow">
      <slot />
    </div>
  </div>

  <footer>
    <div class="text-center font-small drop-shadow-md">
      2012 - 2024
    </div>
  </footer>
</div>
