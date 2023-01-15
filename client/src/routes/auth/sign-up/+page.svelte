<script lang="ts">
	import axios from 'axios';
	import { onMount } from 'svelte';
	let name = '';
	let password = '';
	let email = '';
	let loading = false;

	async function getMe() {
		try {
			console.log('hello bitch');
			const data = await axios.get('http://localhost:4000/me', {
				withCredentials: true
			});
			console.log(data);
		} catch (error) {
			console.log('did not make it :(');
			console.log(error);
		}
	}

	onMount(getMe);

	async function signUp() {
		loading = true;
		const res = await axios.post(
			'http://localhost:4000/sign-up',
			{
				name,
				password,
				email
			},
			{
				withCredentials: true,
				headers: { 'Content-Type': 'text/plain' }
			}
		);

		console.log('gotttiiitt');

		console.log(res);

		loading = false;

		getMe();
	}
</script>

<h1 class="text-5xl">sign up page</h1>
<div class="my-5">
	<h2 class="mb-2 text-2xl font-semibold md:text-3xl">What is your name?</h2>
	<input
		type="text"
		placeholder="adam adamson..."
		bind:value={name}
		class="border-text-dimmed-dark  focus:border-secondary bg-bg mb-10 border-2 py-2 px-5 outline-none"
	/>
	<h2 class="mb-2 text-2xl font-semibold md:text-3xl">What is your email?</h2>
	<input
		type="email"
		placeholder="adam@adam.se..."
		bind:value={email}
		class="border-text-dimmed-dark  focus:border-secondary bg-bg mb-10 border-2 py-2 px-5 outline-none"
	/>

	<h2 class="mb-2 text-2xl font-semibold md:text-3xl">Give us super secure password!</h2>
	<input
		type="text"
		placeholder="password..."
		bind:value={password}
		class="border-text-dimmed-dark focus:border-secondary bg-bg border-2 py-2 px-5 outline-none"
	/>
</div>
<button
	on:click={signUp}
	disabled={!name || !location || loading}
	class="bg-primary flex items-center gap-2 py-2 px-6 transition-opacity disabled:opacity-50 md:text-xl"
>
	{#if loading}
		<div
			class="border-t-secondary h-4 w-4 flex-shrink-0 animate-spin rounded-full border-4 border-black"
		/>
	{/if}
	Create Stash
</button>
