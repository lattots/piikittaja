import { User } from "./models";

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

