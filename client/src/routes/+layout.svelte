<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import Cart from '$lib/components/cart/Cart.svelte';
	import '../app.css';
	import axios from 'axios';

	// axios.defaults.baseURL = 'http://localhost:4000';
</script>

<div class="px-4">
	<div class="mx-auto max-w-6xl">
		<Navbar />
		<div class="py-10 md:py-14">
			<slot />
		</div>
		<Cart />
	</div>
</div>
