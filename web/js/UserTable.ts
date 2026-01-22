import { fetchUsers } from "./api"
import { User } from "./models"
import { format } from "./monetaryUtil"
import "./UserModal"

export class UserTable extends HTMLElement {
	apiUrl: string = ""
	private tbody: HTMLElement | null = null;

	private handleSearchBound = this.handleSearchSuccess.bind(this);

	constructor() {
		super()
	}

	async connectedCallback() {
		this.apiUrl = this.getAttribute("api-url") || ""

		this.renderSkeleton();

		await this.loadInitialUsers();

		window.addEventListener("search-success", this.handleSearchBound as EventListener);
	}

	disconnectedCallback() {
		window.removeEventListener("search-success", this.handleSearchBound as EventListener);
	}

	renderSkeleton() {
		this.innerHTML = `
			<table id="user-table">
				<thead style="font-weight: bold">
					<tr>
						<th>Nimi</th>
						<th>Telegram</th>
						<th>Saldo</th>
					</tr>
				</thead>
				<tbody id="user-table-body">
					<tr><td colspan="3">Ladataan...</td></tr>
				</tbody>
			</table>
			<user-modal id="global-modal" api-url="${this.apiUrl}"></user-modal>
		`;

		this.tbody = this.querySelector("#user-table-body");
		this.setupEventListeners();
	}

	async loadInitialUsers() {
		try {
			const users = await fetchUsers(this.apiUrl);
			this.updateRows(users);
		} catch (error) {
			console.error("Failed to load initial users:", error);
			if (this.tbody) this.tbody.innerHTML = `<tr><td colspan="3">Käyttäjien lataaminen epäonnistui</td></tr>`;
		}
	}

	handleSearchSuccess(e: CustomEvent) {
		const users = e.detail.users;
		this.updateRows(users);
	}

	updateRows(users: User[]) {
		if (!this.tbody) return;

		if (users.length === 0) {
			this.tbody.innerHTML = `<tr><td colspan="3">Ei käyttäjiä</td></tr>`;
			return;
		}

		this.tbody.innerHTML = renderUsers(users);
	}

	setupEventListeners() {
		const tbody = this.querySelector("tbody");
		if (!tbody) return;

		tbody.addEventListener("click", (e) => {
			const target = e.target as Element;
			const row = target.closest(".user-row") as HTMLTableRowElement | null;
			if (!row) return;

			const userId = row.getAttribute("data-id");
			const modal = this.querySelector("#global-modal") as any;

			if (modal && userId) {
				modal.setAttribute("user-id", userId);
				modal.open();
			}
		});
	}
}

export function renderUsers(users: User[]): string {
	return users.map(user => `
		<tr class="user-row" data-id="${user.id}" style="cursor: pointer;">
			<td>${user.firstName} ${user.lastName}</td>
			<td>${user.username}</td>
			<td style="text-align: right; padding-right: 1.2rem">${format(user.balance)}</td>
		</tr>
	`).join("")
}

customElements.define("user-table", UserTable)
