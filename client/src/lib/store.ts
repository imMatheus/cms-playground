import { writable } from 'svelte/store';
import type { Product } from '$lib/types';

interface CartStore {
	isOpen: boolean;
	data: Product[];
}

export const cartStore = writable<CartStore>({ isOpen: false, data: [] });
