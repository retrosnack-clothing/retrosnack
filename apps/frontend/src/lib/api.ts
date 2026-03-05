const BASE_URL = import.meta.env.PUBLIC_API_URL ?? 'http://localhost:8080';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		headers: { 'Content-Type': 'application/json', ...init?.headers },
		...init
	});
	if (!res.ok) {
		const body = await res.text();
		throw new Error(`HTTP ${res.status}: ${body}`);
	}
	return res.json() as Promise<T>;
}

export const api = {
	products: {
		list: () => request<Product[]>('/api/products'),
		get: (id: string) => request<Product>(`/api/products/${id}`)
	},
	categories: {
		list: () => request<Category[]>('/api/categories')
	}
};

// Types — kept in sync with Go models
export interface Product {
	id: string;
	title: string;
	description: string;
	brand: string;
	condition: 'excellent' | 'good' | 'fair';
	price_cents: number;
	category_id: string;
	instagram_post_url: string;
	images: ProductImage[];
	created_at: string;
}

export interface ProductImage {
	id: string;
	url: string;
	position: number;
}

export interface Category {
	id: string;
	name: string;
	slug: string;
	parent_id: string | null;
}
