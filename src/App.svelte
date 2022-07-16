<script>
	import MediaBox from "./MediaBox.svelte";
	import ImageGrid from "./ImageGrid.svelte";

	import { getImages } from "./api.js";

	const images = getImages();
</script>

<main>
	<ImageGrid images={images} columns=4/>
</main>

<style>
</style>
