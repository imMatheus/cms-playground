<script lang="ts">
	import type { Product } from '$lib/types';
	import ProductCard from '$lib/components/ProductCard.svelte';

	export let data: { props: Product[] };
	const products = data.props;
	console.log(data);

	// onMount(async () => {
	// 	console.log('wag 1');

	// 	const res = await fetch('http://localhost:4000/stashes');
	// 	const d = await res.json();
	// 	console.log(d);
	// });
</script>

<h1 class="mb-6 text-7xl font-semibold">Find good products,<br /> near you</h1>

<div
	class="grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:gap-x-8"
>
	{#each products as product}
		<ProductCard {product} />
	{/each}
</div>

<pre>{JSON.stringify(products, null, 2)}</pre>
