<script>
	import MediaBox from "./MediaBox.svelte";
	import { writable } from "svelte/store";

	const pepe = "https://image.cnbcfm.com/api/v1/image/103978904-The_meme_formerly_known_as_Kuk_1.png?v=1475149183&w=740&h=416";

	let y = 0;
	const columns = 4;
	const pepes = 1000;
	const rows = Math.ceil(pepes / columns);

	let minId = 0;
	let maxId = 0;

	$: {
		// This doesn't have to be the window height, it could be a component.
		const totalScroll = document.documentElement.scrollHeight - window.innerHeight;
		if (Number.isFinite(y / totalScroll)) {
			let minRow = Math.floor(y / document.documentElement.scrollHeight * rows);
			let maxRow = Math.ceil((y + window.innerHeight) / document.documentElement.scrollHeight * rows);
			minId = (minRow + 1) * columns
			maxId = (maxRow - 1) * columns
		}
	}
</script>

<svelte:window bind:scrollY={y}/>

<main>
	<div id="grid">
	{#each Array.from(Array(pepes).keys()) as i}
		<MediaBox src={pepe} render={(minId <= i) && (i < maxId)}/>
	{/each}
</main>

<style>
#grid {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	padding: 0px 30%;
	background-color: red;
}
</style>
