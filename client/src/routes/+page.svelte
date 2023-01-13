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

<h1 class="mb-6 font-serif text-7xl font-semibold">Find good products,<br /> near you</h1>

<div class="grid grid-cols-4 md:gap-10 gap-2 bg-secondary">
	{#each products as product}
		<a href="/product/abc">
			<div class="bg-primary p-3">{product.name}</div>
		</a>
	{/each}
</div>

<pre>{JSON.stringify(products, null, 2)}</pre>
