import type { Load } from '@sveltejs/kit';
import type { Product } from '$lib/types';
import axios from 'axios';

export const load: Load = async (res) => {
	const data = await axios.get('http://localhost:4000/products');
	console.log(data.data);

	return {
		props: data.data
	};
};
