import { User } from "./models.ts";
import { format } from "./monetaryUtil.ts";
import "./TransactionTable.ts";
import "./TransactionForm.ts";

import userStyles from "../css/user.css";
import generalStyles from "../css/general.css";

const userSheet = new CSSStyleSheet();
userSheet.replaceSync(userStyles);

const generalSheet = new CSSStyleSheet();
generalSheet.replaceSync(generalStyles);

export class UserModal extends HTMLElement {
  private dialog: HTMLDialogElement | null = null;

  private shadow: ShadowRoot;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });

    this.shadow.adoptedStyleSheets = [generalSheet, userSheet];
  }

  connectedCallback() {
    this.render();
    this.dialog = this.shadow.querySelector("dialog") || null;

    if (this.dialog) {
      this.dialog.addEventListener("click", (e) => {
        if (!this.dialog) return;

        const rect = this.dialog.getBoundingClientRect();

        const isClickOutside = e.clientX < rect.left ||
          e.clientX > rect.right ||
          e.clientY < rect.top ||
          e.clientY > rect.bottom;

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

    const content = this.shadow.querySelector("#content");
    if (content) content.innerHTML = "<p>Loading details...</p>";

    this.dialog.showModal();

    try {
      const resp = await fetch(`${apiUrl}/users/${userId}`);
      if (resp.status === 401 || resp.status === 403) {
        window.location.href = "/login";
        return;
      }
      const user = await resp.json();
      this.updateContent(user);
    } catch (err) {
      if (content) content.innerHTML = "<p>Error loading user.</p>";
    }
  }

  private updateContent(user: User) {
    const apiUrl = this.getAttribute("api-url");
    const content = this.shadow.querySelector("#content");
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
    if (!this.shadow) return;
    this.shadow.innerHTML = `
         <dialog class="modal">
              <div id="content"></div>
          </dialog>
        `;
  }
}

function renderUserInfo(user: User): string {
  const name: string = (user.firstName)
    ? `${user.firstName} ${user.lastName}`
    : "Matti Meikalainen";
  const balanceColorClass: string = (user.balance >= 0) ? "green" : "red";
  console.log(balanceColorClass);
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
