<script>
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { loc, globalMap, globalChallenge, globalResult } from './stores.js';

	import CreateMap from './CreateMap.svelte'
	import CreateChallenge from './CreateChallenge.svelte'
	import Resume from './Resume.svelte'
	import Join from './Join.svelte'
	import Scores from './Scores.svelte'
	import Summary from './Summary.svelte'
	// TODO: code split this out into a separate bundle
	import Modify from './Modify.svelte'

	let ewapi = new EarthwalkerAPI();

	// TODO: remove debug
	$: console.log($loc);
	$: console.log($globalMap);
	$: console.log($globalChallenge);
	$: console.log($globalResult);

	// write sets loc without side effects (unlike set/assignment)
	loc.write(window.location.pathname);
</script>

<svelte:window on:popstate={(e) => $loc = e.target.location.pathname} />

<style>
	#content {
		margin: 2em;
	}
</style>

<main>
	{#if $loc.startsWith("/play")}
		<!-- TODO: code split this out into a separate bundle -->
		<Modify {ewapi}/>
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
		<div id="content">
			{#if $loc === "/"}
				<Resume/>
			{:else if $loc.startsWith("/createmap")}
				<CreateMap/>
			{:else if $loc.startsWith("/createchallenge")}
				<CreateChallenge {ewapi}/>
			{:else if $loc.startsWith("/join")}
				<Join {ewapi}/>
			{:else if $loc.startsWith("/scores")}
				<Scores {ewapi}/>
			{:else if $loc.startsWith("/summary")}
				<Summary {ewapi}/>
			{:else}
				<h3>404.  That's an error.</h3>
			{/if}
		</div>
	{/if}
</main>