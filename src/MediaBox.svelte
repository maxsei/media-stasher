<script>
	export let render = false;
	export let src = "";
</script>

<div id="container">
	{#if render}
		<img src={src}/>
	{/if}
</div>

<style>
#container{
	display: flex;
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
