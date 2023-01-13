<script lang="ts">
	let name = '';
	let location = '';
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
		disabled={!name || !location}
		class="bg-primary py-2 px-6 transition-opacity disabled:opacity-50 md:text-xl"
	>
		Create Stash
	</button>
</div>
