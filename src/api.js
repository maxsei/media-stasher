export function getImages(){
	const pepe = "https://image.cnbcfm.com/api/v1/image/103978904-The_meme_formerly_known_as_Kuk_1.png?v=1475149183&w=740&h=416";
	const n = 1000;
	return Array.from(Array(n).keys()).map(()=>({src: pepe}));
}
