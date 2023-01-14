<script lang="ts">
	import axios from 'axios';
	let name = '';
	let location = '';
	let loading = false;

	async function createStash() {
		loading = true;
		const res = await axios.post(
			'http://localhost:4000/stash',
			{
				name,
				location
			},
			{
				headers: { 'Content-Type': 'text/plain' }
			}
		);

		console.log('gotttiiitt');

		console.log(res);

		loading = false;
	}
</script>

<div class="">
	<h1 class="mb-10 text-4xl font-bold md:mb-14 md:text-6xl">Let's create a Stash!</h1>

	<div class="my-5">
		<h2 class="mb-2 text-2xl font-semibold md:text-3xl">What should we call your Stash?</h2>
		<input
			type="text"
			placeholder="my stash..."
			bind:value={name}
			class="border-text-dimmed-dark  focus:border-secondary bg-bg mb-10 border-2 py-2 px-5 outline-none"
		/>

		<h2 class="mb-2 text-2xl font-semibold md:text-3xl">Where is your Stash located?</h2>
		<input
			type="text"
			placeholder="stockholm sweden..."
			bind:value={location}
			class="border-text-dimmed-dark focus:border-secondary bg-bg border-2 py-2 px-5 outline-none"
		/>
	</div>
	<button
		on:click={createStash}
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
</div>
