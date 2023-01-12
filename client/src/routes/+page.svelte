<script lang="ts">
	import { onMount } from 'svelte';
	import type { Product } from '$lib/types';

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

<div class="px-4 py-10 md:py-14">
	<div class="max-w-6xl mx-auto">
		<h1 class="text-7xl font-semibold font-serif">Find good products,<br /> near you</h1>

		<div class="grid bg-secondary grid-cols-4">
			{#each products as product}
				<div class="p-3 bg-primary">{product.name}</div>
			{/each}
		</div>

		<pre>{JSON.stringify(products, null, 2)}</pre>
	</div>
</div>
