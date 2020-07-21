<script>
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { loc, ewapi, globalMap, globalChallenge, globalResult } from './stores.js';

	import CreateMap from './CreateMap.svelte'
	import CreateChallenge from './CreateChallenge.svelte'
	import Resume from './Resume.svelte'
	import Join from './Join.svelte'
	import Scores from './Scores.svelte'
	import Summary from './Summary.svelte'
	// TODO: code split this out into a separate bundle
	import Modify from './Modify.svelte'

	$ewapi = new EarthwalkerAPI();

    // TODO: FIXME: a cleaner way with no race conditions.  
    //       Derived stores with promises/callbacks?
    async function setResultChallengeMap(resultID) {
        $globalResult = await $ewapi.getResult(resultID);
        if (!$globalChallenge || $globalResult.ChallengeID !== $globalChallenge.ChallengeID) {
            return setChallengeMap($globalResult.ChallengeID);
        }
    }

    async function setChallengeMap(challengeID) {
        $globalChallenge = await $ewapi.getChallenge(challengeID);
        if (!$globalMap || $globalChallenge.MapID !== $globalMap.MapID) {
            $globalMap = await $ewapi.getMap($globalChallenge.MapID);
        }
    }

    onMount(async () => {
		let challengeID = getChallengeID();
		if (challengeID) {
			await setResultChallengeMap(getChallengeResultID(challengeID));
		}
    });

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
		{#if $globalMap && $globalChallenge && $globalResult}
			<!-- TODO: code split this out into a separate bundle -->
			<Modify/>
		{:else}
			<h3>Loading...</h3>
		{/if}
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
			{#if $globalMap && $globalChallenge && $globalResult}
				{#if $loc === "/"}
					<Resume/>
				{:else if $loc.startsWith("/createmap")}
					<CreateMap/>
				{:else if $loc.startsWith("/createchallenge")}
					<CreateChallenge/>
				{:else if $loc.startsWith("/join")}
					<Join/>
				{:else if $loc.startsWith("/scores")}
					<Scores/>
				{:else if $loc.startsWith("/summary")}
					<Summary/>
				{:else}
					<h3>404.  That's an error.</h3>
				{/if}
			{:else}
				<h3>Loading...</h3>
			{/if}
		</div>
	{/if}
</main>