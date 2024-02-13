<script lang="ts">
	import '$lib/styles.scss';
</script>

<div class="app">
	<main>
		<slot />
	</main>

	<footer>
	</footer>
</div>

<style>



main {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}
</style>
