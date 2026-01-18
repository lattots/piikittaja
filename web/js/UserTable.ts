import { User } from "./models"
import { format } from "./monetaryUtil"
import "./UserModal"

export class UserTable extends HTMLElement {
	apiUrl: string = ""

	constructor() {
		super()
	}

	async connectedCallback() {
		this.apiUrl = this.getAttribute("api-url") || ""

		this.innerHTML = `<p>Loading users...</p>`
		await this.render()
	}

	async render() {
		const resp = await fetch(`${this.apiUrl}/users`)
		if (!resp.ok) return;

		const users: User[] = await resp.json()

		this.innerHTML = `
            <table id="user-table">
                <thead style="font-weight: bold">
                    <tr>
						<th>Nimi</th>
						<th>Telegram</th>
						<th>Saldo</th>
					</tr>
                </thead>
                <tbody>
                    ${users.map(user => `
                        <tr class="user-row" data-id="${user.id}" style="cursor: pointer;">
                            <td>${user.firstName} ${user.lastName}</td>
                            <td>${user.username}</td>
                            <td style="text-align: right; padding-right: 8px">${format(user.balance)}</td>
                        </tr>
                    `).join("")}
                </tbody>
            </table>
            <user-modal id="global-modal" api-url="${this.apiUrl}"></user-modal>
        `;

		this.setupEventListeners();
	}

	setupEventListeners() {
		const tbody = this.querySelector('tbody');
		if (!tbody) return;

		tbody.addEventListener('click', (e) => {
			const row = (e.target as HTMLElement).closest('.user-row');
			if (!row) return;

			const userId = row.getAttribute('data-id');
			const modal = this.querySelector('#global-modal') as any;

			if (modal && userId) {
				modal.setAttribute('user-id', userId);
				modal.open();
			}
		});
	}
}

customElements.define("user-table", UserTable)
