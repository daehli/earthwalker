<script>
	import CreateMap from './CreateMap.svelte'
	import CreateChallenge from './CreateChallenge.svelte'
	$: console.log(window.location.pathname)
</script>

<main>
	<nav class="navbar navbar-expand-sm navbar-light bg-light">
		<span class="navbar-brand">Earthwalker</span>
		<ul class="navbar-nav">
			<div class="collapse navbar-collapse">
				<li class="nav-item active">
					<a class="nav-link" href="/">Home</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" href="https://gitlab.com/glatteis/earthwalker">Source code</a>
				</li>
			</div>
		</ul>
	</nav>
	{#if window.location.pathname === "/"}
		<p>Landing page.</p>
	{:else if window.location.pathname === "/createmap"}
		<CreateMap/>
	{:else if window.location.pathname === "/createchallenge"}
		<CreateChallenge/>
	{:else}
		<h3>404.  That's an error.</h3>
	{/if}
</main>

<style>
	main {
		margin: 0;
	}
</style>