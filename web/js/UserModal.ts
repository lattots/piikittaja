import { User } from "./models";
import { format } from "./monetaryUtil";
import "./TransactionTable"
import "./TransactionForm"

export class UserModal extends HTMLElement {
	private dialog: HTMLDialogElement | null = null;

	constructor() {
		super();
		this.attachShadow({ mode: "open" });
	}

	connectedCallback() {
		this.render();
		this.dialog = this.shadowRoot?.querySelector("dialog") || null;

		if (this.dialog) {
			this.dialog.addEventListener('click', (e) => {
				if (!this.dialog) return;

				const rect = this.dialog.getBoundingClientRect();

				const isClickOutside = (
					e.clientX < rect.left ||
					e.clientX > rect.right ||
					e.clientY < rect.top ||
					e.clientY > rect.bottom
				);

				if (isClickOutside) {
					this.dialog.close();
				}
			});
		}
	}

	async open() {
		if (!this.dialog) return;

		const userId = this.getAttribute("user-id");
		const apiUrl = this.getAttribute("api-url");

		const content = this.shadowRoot?.querySelector("#content");
		if (content) content.innerHTML = "<p>Loading details...</p>";

		this.dialog.showModal();

		try {
			const resp = await fetch(`${apiUrl}/users/${userId}`);
			const user = await resp.json();
			this.updateContent(user);
		} catch (err) {
			if (content) content.innerHTML = "<p>Error loading user.</p>";
		}
	}

	private updateContent(user: User) {
		const apiUrl = this.getAttribute("api-url");
		const content = this.shadowRoot?.querySelector("#content");
		if (!content) return;

		content.innerHTML = `
			${renderUserInfo(user)}
			<transaction-form api-url=${apiUrl} user-id=${user.id}></transaction-form>
			<transaction-table api-url=${apiUrl} user-id=${user.id}></transaction-table>
        `;

		const txForm = content.querySelector("transaction-form");
		txForm?.addEventListener("transaction-success", (e: Event) => {
			const customEvent = e as CustomEvent<User>;
			const updatedUser = customEvent.detail;

			console.log("New balance received:", updatedUser.balance);

			this.updateContent(updatedUser);
		});
	}

	private render() {
		if (!this.shadowRoot) return;
		this.shadowRoot.innerHTML = `
		<style>
			@import url("style.css");
		</style>
           <dialog class="modal">
                <div id="content"></div>
            </dialog>
        `;
	}
}

function renderUserInfo(user: User): string {
	const name: string = (user.firstName) ? `${user.firstName} ${user.lastName}` : "Matti Meikalainen";
	const balanceColorClass: string = (user.balance >= 0) ? "green" : "red";
	return `
		<h2>${user.username}</h2>
		<div class="flexbox" style="justify-content: space-between">
			<p style="font-weight: bold; align-content: center">${name}</p>
			<div class="${balanceColorClass}" id="balance-display">
				<p style="font-weight: bold">${format(user.balance)}</p>
			</div>
		</div>
	`;
}

customElements.define("user-modal", UserModal);
