import type { Load } from '@sveltejs/kit';
import type { Product } from '$lib/types';
import axios from 'axios';
import { error } from '@sveltejs/kit';

export const load: Load = async ({ params }) => {
	if (params.slug) {
		const data = await axios.get('http://localhost:4000/products/' + params.slug);
		console.log(data.data);

		return {
			props: data.data
		};
	}

	throw error(404, 'Not found');
};
