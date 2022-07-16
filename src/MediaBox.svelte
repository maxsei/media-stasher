<script>
	export let containerBounds;
	export let src = "";

	let component;
	// Check if image should be rendered based on the container bounds.
	let render = false;
	$: {
		const top = component?.offsetTop ?? 0;
		const bottom = top + component?.offsetHeight ?? 0;
		const containerBoundsTop = containerBounds?.top ?? 0;
		const containerBoundsBottom = containerBounds?.bottom ?? 0;
		render = (
			(component !== undefined) 
			&& (containerBounds !== undefined) 
			&& (
				(containerBoundsTop <= bottom)
				&& (top <= containerBoundsBottom)
			)
		);
	}
</script>

<div id="container" bind:this={component}>
	{#if render}
		<img src={src}/>
	{/if}
</div>

<style>
#container{
	display: flex;
	position: relative;
	justify-content: center;
	width: 100px;
	height: 100px;
	background-color: blue;
}
img {
	object-fit: scale-down;
	max-width: 100%;
	max-height: 100%;
}
</style>
