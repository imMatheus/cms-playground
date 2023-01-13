<script lang="ts">
	import type { Product } from '$lib/types';

	export let product: Product;
</script>

<a href={`/product/${product.id}`}>
	<div class="bg-bg group relative">
		<div class="w-full overflow-hidden rounded-md bg-gray-200 group-hover:opacity-80 lg:h-80">
			<img
				src="https://tailwindui.com/img/ecommerce-images/product-page-01-related-product-01.jpg"
				alt="Front of men&#039;s Basic Tee in black."
				class="h-full w-full object-cover object-center lg:h-full lg:w-full"
			/>
		</div>
		<div class="mt-4 flex justify-between">
			<div>
				<h3 class="text-sm text-gray-700">
					<span aria-hidden="true" class="absolute inset-0" />
					{product.name}
				</h3>
				<p class="mt-1 text-sm text-gray-500">Mendes store</p>
			</div>
			<p class="text-sm font-medium text-gray-900">${product.price}</p>
		</div>
	</div>
</a>
