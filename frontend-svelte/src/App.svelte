<script>
	import CreateMap from './CreateMap.svelte'
	import CreateChallenge from './CreateChallenge.svelte'
	import Resume from './Resume.svelte'
	import Join from './Join.svelte'
	import Scores from './Scores.svelte'

	// TODO: code split this out into a separate bundle
	import Modify from './Modify.svelte'
</script>

<style>
</style>

<main>
	{#if window.location.pathname.startsWith("/play")}
		<!-- TODO: code split this out into a separate bundle -->
		<Modify/>
	{:else}
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
			<Resume/>
		{:else if window.location.pathname === "/createmap"}
			<CreateMap/>
		{:else if window.location.pathname === "/createchallenge"}
			<CreateChallenge/>
		{:else if window.location.pathname === "/join"}
			<Join/>
		{:else if window.location.pathname === "/scores"}
			<Scores/>
		{:else}
			<h3>404.  That's an error.</h3>
		{/if}
	{/if}
</main>