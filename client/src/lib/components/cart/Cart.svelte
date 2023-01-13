<script language="ts">
	import { cartStore } from '$lib/store';
	import CartItem from './CartItem.svelte';
	import classNames from 'classnames';
</script>

<div class="relative z-10" aria-labelledby="slide-over-title" role="dialog" aria-modal="true">
	<!-- <div
		class={classNames(
			'fixed inset-0 bg-red-500 bg-opacity-75 transition-opacity duration-500',
			!$cartStore.isOpen && 'bg-opacity-20'
		)}
	/> -->

	<!-- <div class="fixed inset-0 overflow-hidden">
		<div class="absolute inset-0 overflow-hidden"> -->
	<div
		class={classNames(
			'pointer-events-none fixed inset-y-0 right-0 flex max-w-full transform pl-10 transition duration-500 ease-in-out sm:duration-700',
			!$cartStore.isOpen && 'translate-x-full'
		)}
	>
		<div class="pointer-events-auto w-screen max-w-md">
			<div class="flex h-full flex-col overflow-y-scroll bg-white shadow-xl">
				<div class="flex-1 overflow-y-auto py-6 px-4 sm:px-6">
					<div class="flex items-start justify-between">
						<h2 class="text-lg font-medium text-gray-900" id="slide-over-title">Shopping cart</h2>
						<div class="ml-3 flex h-7 items-center">
							<button
								type="button"
								on:click={() => ($cartStore.isOpen = false)}
								class="-m-2 p-2 text-gray-400 hover:text-gray-500"
							>
								<span class="sr-only">Close panel</span>
								<!-- Heroicon name: outline/x-mark -->
								<div class="h-6 w-6 bg-red-500" />
							</button>
						</div>
					</div>

					<div class="mt-8">
						<div class="flow-root">
							<ul class="-my-6 divide-y divide-gray-200">
								<CartItem />
								<CartItem />
								<CartItem />
								<CartItem />
								<CartItem />
								<CartItem />
								<CartItem />
							</ul>
						</div>
					</div>
				</div>

				<div class="border-t border-gray-200 py-6 px-4 sm:px-6">
					<div class="flex justify-between text-base font-medium text-gray-900">
						<p>Subtotal</p>
						<p>$262.00</p>
					</div>
					<p class="mt-0.5 text-sm text-gray-500">Shipping and taxes calculated at checkout.</p>
					<div class="mt-6">
						<a
							href="#"
							class="flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-6 py-3 text-base font-medium text-white shadow-sm hover:bg-indigo-700"
							>Checkout</a
						>
					</div>
					<div class="mt-6 flex justify-center text-center text-sm text-gray-500">
						<p>
							or
							<button type="button" class="font-medium text-indigo-600 hover:text-indigo-500">
								Continue Shopping
								<span aria-hidden="true"> &rarr;</span>
							</button>
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
	<!-- </div>
	</div> -->
</div>
