<script>
	import MediaBox from "./MediaBox.svelte";


	export let images;
	export let columns;

	let component;
	let containerBounds;
	function setContainerBounds(){
		const top = component?.scrollTop ?? 0;
		const bottom = top + component?.clientHeight ?? 0;
		containerBounds = {top, bottom};
	}
	$: if (component !== undefined) { setContainerBounds(); }
</script>

<div 
	id="grid"
	style="grid-template-columns: repeat({columns}, 1fr);"
	on:scroll={setContainerBounds}
	bind:this={component}
>
	{#each images as image}
		<MediaBox {...{...image, containerBounds}}/>
	{/each}
</div>

<style>
#grid {
	position: relative;
	display: grid;
	padding: 0px 30%;
	background-color: red;
	overflow-y: scroll;
	height: 600px;
}
</style>
