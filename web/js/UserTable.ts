import { User } from "./models"

export class UserTable extends HTMLElement {
	apiUrl: string = ""

	constructor() {
		super()
	}

	async connectedCallback() {
		this.apiUrl = this.getAttribute('api-url') || ""

		this.innerHTML = `<p>Loading users...</p>`
		await this.render()
	}

	async render() {
		const resp = await fetch(`${this.apiUrl}/users`)
		if (!resp.ok) {
			// TODO: Handle error
			return
		}
		const users: User[] = await resp.json()

		this.innerHTML = `
      <table>
        <thead>
          <tr>
            <th>Telegram</th>
            <th>Saldo</th>
          </tr>
        </thead>
        <tbody>
          ${users.map(user => `
            <tr>
              <td>${user.username}</td>
              <td>${balanceString(user.balance)}</td>
            </tr>
          `).join("")}
        </tbody>
      </table>
    `;
	}
}

function balanceString(balance: number): string {
	const balanceEuros: number = balance / 100
	return `${balanceEuros} €`
}

customElements.define('user-table', UserTable)
