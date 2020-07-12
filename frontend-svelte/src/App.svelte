<script>
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { loc } from './stores.js';

	import CreateMap from './CreateMap.svelte'
	import CreateChallenge from './CreateChallenge.svelte'
	import Resume from './Resume.svelte'
	import Join from './Join.svelte'
	import Scores from './Scores.svelte'
	import Summary from './Summary.svelte'
	// TODO: code split this out into a separate bundle
	import Modify from './Modify.svelte'

	let ewapi = new EarthwalkerAPI();
	let challengeID = getChallengeID();
	let resultID = getChallengeResultID(challengeID);
	let curChallenge;
	$: if (challengeID) {
		ewapi.getChallenge(challengeID).then(challenge => curChallenge = challenge || null);
	};
	let curMap;
	$: if (curChallenge && curChallenge.MapID) {
		ewapi.getMap(curChallenge.MapID).then(map => curMap = map || null);
	};
	let curResult;
	$: if (resultID) {
		ewapi.getResult(resultID).then(result => curResult = result || null);
	};

	// TODO: remove debug
	$: console.log(curMap);
	$: console.log(curResult);
	$: console.log(curChallenge);

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
	{#if $loc.startsWith("/play") && ewapi && curMap && curChallenge && curResult}
		<!-- TODO: code split this out into a separate bundle -->
		<Modify {ewapi} {curMap} {curChallenge} {curResult}/>
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
			{#if $loc === "/" && curChallenge && curResult}
				<Resume {curChallenge} {curResult}/>
			{:else if $loc.startsWith("/createmap")}
				<CreateMap/>
			{:else if $loc.startsWith("/createchallenge") && ewapi && curMap && curChallenge}
				<CreateChallenge {ewapi} {curMap} {curChallenge}/>
			{:else if $loc.startsWith("/join") && ewapi && curChallenge && curResult}
				<Join {ewapi} {curChallenge} {curResult}/>
			{:else if $loc.startsWith("/scores") && ewapi && curMap && curChallenge && curResult}
				<Scores {ewapi} {curMap} {curChallenge} {curResult}/>
			{:else if $loc.startsWith("/summary") && ewapi && curMap && curChallenge && curResult}
				<Summary {ewapi} {curMap} {curChallenge} {curResult}/>
			{:else}
				<h3>404.  That's an error.</h3>
			{/if}
		</div>
	{/if}
</main>