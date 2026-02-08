import { ApiError } from "./errors";
import { Transaction, User } from "./models";

export async function fetchUsers(apiUrl: string, searchTerm?: string, signal?: AbortSignal): Promise<User[]> {
	const url = new URL(`${apiUrl}/users`, window.location.origin);

	if (searchTerm) {
		url.searchParams.set("searchTerm", searchTerm);
	}

	const resp = await fetch(url.href, { signal });

	if (resp.status === 401 || resp.status === 403) {
		window.location.href = '/login';
		return [];
	}

	if (!resp.ok) return [];

	const users: User[] = await resp.json();
	return users;
}

export async function fetchUserTransactions(apiUrl: string, userId: string, quantity?: number): Promise<Transaction[]> {
	const resp = await fetch(`${apiUrl}/users/${userId}/transactions?quantity=${quantity}`)
	if (resp.status === 401) {
		throw new ApiError("UNAUTHENTICATED_ERROR", "user not authenticated");
	}
	if (resp.status === 403) {
		throw new ApiError("FORBIDDEN_ERROR", "user is not allowed to access this resource")
	}
	if (!resp.ok) throw Error("unexpected error fetching user transactions");

	const rawData: any[] = await resp.json();

	if (!rawData || rawData.length === 0) {
		return [];
	}

	const transactions: Transaction[] = rawData.map(t => ({
		...t,
		issuedAt: new Date(t.issuedAt)
	}));

	return transactions;
}

export async function fetchTransactions(apiUrl: string, endDate: Date, timeWindowDays: number, type?: string): Promise<Transaction[]> {
	const url = new URL(`${apiUrl}/transactions`, window.location.origin);
	url.searchParams.set("endDate", endDate.toISOString());
	url.searchParams.set("window", timeWindowDays.toString());
	if (type) url.searchParams.set("type", type);

	const resp = await fetch(url)
	if (resp.status === 401) {
		throw new ApiError("UNAUTHENTICATED_ERROR", "user not authenticated");
	}
	if (resp.status === 403) {
		throw new ApiError("FORBIDDEN_ERROR", "user is not allowed to access this resource")
	}
	if (!resp.ok) throw Error("unexpected error fetching user transactions");

	const rawData: any[] = await resp.json();

	if (!rawData || rawData.length === 0) {
		return [];
	}

	const transactions: Transaction[] = rawData.map(t => ({
		...t,
		issuedAt: new Date(t.issuedAt)
	}));

	return transactions;
}

