<script>
	import { onMount } from "svelte";

	import MediaBox from "./MediaBox.svelte";
	import ImageGrid from "./ImageGrid.svelte";

	import { getImages } from "./api.js";

	// const images = getImages();
	let images = [];
	onMount(async () =>{
		const response = await fetch("/filepaths");
		const data = await response.json();
		images = data.paths.map((el)=> ({src: `thumbnail/${el}`}));
	})

</script>

<main>
	<ImageGrid images={images} columns=4/>
</main>

<style>
</style>
