<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.css';
	import axios from 'axios';

	// axios.defaults.baseURL = 'http://localhost:4000';
</script>

<Navbar />
<slot />
