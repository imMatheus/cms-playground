<script lang="ts">
	import { ShoppingBagIcon } from 'svelte-feather-icons';
	import { cartStore } from '$lib/store';
</script>

<nav class="pt-5">
	<div class="mx-auto max-w-6xl">
		<div class="flex items-center justify-between">
			<a href="/" class="text-3xl font-semibold">Stash</a>

			<div class="flex gap-4">
				<a href="/stashes" class="transition-opacity hover:opacity-75">Stashes</a>
				<a href="/stashes/mine" class="transition-opacity hover:opacity-75">My stashes</a>
				<button
					class="flex cursor-pointer items-center gap-1 transition-opacity hover:opacity-75"
					on:click={() => {
						$cartStore.isOpen = true;
					}}
				>
					<p>Cart</p>
					<ShoppingBagIcon size="16" />
				</button>
				<a href="/auth/sign-up" class="transition-opacity hover:opacity-75">Sign up</a>
			</div>
		</div>
	</div>
</nav>
