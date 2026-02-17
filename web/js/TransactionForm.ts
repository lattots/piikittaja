import { User } from "./models.ts";

export class TransactionForm extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.innerHTML =
      `<p style="font-weight: bold">Loading transaction form...</p>`;
    this.render();
  }

  static get observedAttributes() {
    return ["api-url", "user-id"];
  }

  render() {
    this.innerHTML = `
			<div class="transaction-container">
				<form id="tx-form" autocomplete="off">
					<div class="input-group">
						<input id="amount-input" type="number" name="amount" min="0" step="0.01" required placeholder="0.00">
						<label for="amount-input">€</label>
					</div>
					<div class="action-buttons">
						<button class="button green" type="submit" name="type" value="deposit">Talleta</button>
						<button class="button red" type="submit" name="type" value="withdraw">Nosta</button>
					</div>
				</form>
				<div id="status"></div>
			</div>
		`;

    const txFormEl = this.querySelector("#tx-form");
    if (!txFormEl) {
      console.log("No tx-form element");
      this.innerHTML = `<p>Something went wrong:(`;
      return;
    }
    txFormEl.addEventListener("submit", (e: Event) => this.handleSubmit(e));
  }

  async handleSubmit(e: Event) {
    if (!e) return;
    e.preventDefault();

    const submitEvent = e as SubmitEvent;
    const submitter = submitEvent.submitter as HTMLButtonElement;
    const type = submitter.value; // "deposit" or "withdraw"

    const form = e.target as HTMLFormElement;
    if (!form) return;

    const statusEl = this.querySelector("#status") as HTMLElement;
    if (!statusEl) return;

    const formData = new FormData(form);
    const amountStr = formData.get("amount") as string;
    if (!amountStr) {
      console.log("No amount in form");
      statusEl.innerHTML = "<p>Something went wrong:(";
      return;
    }
    const amount = parseFloat(amountStr);
    const amountInCents = Math.round(amount * 100);

    statusEl.textContent = "Sending...";
    statusEl.style.color = "#666";

    const apiUrl = this.getAttribute("api-url");
    const userId = this.getAttribute("user-id");

    try {
      const response = await fetch(`${apiUrl}/users/${userId}/transactions`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          amount: amountInCents,
          type: type,
        }),
      });

      if (response.status === 401 || response.status === 403) {
        window.location.href = "/login";
        return;
      }

      if (response.ok) {
        form.reset();

        const updatedUser: User = await response.json();

        this.dispatchEvent(
          new CustomEvent("transaction-success", {
            bubbles: true,
            composed: true,
            detail: updatedUser,
          }),
        );
      } else if (response.status === 402) {
        statusEl.textContent = "Rahat on loppu:(";
        statusEl.style.color = "#FF8270";
      } else {
        const errorMsg: string = await response.text();
        statusEl.textContent = `Tapahtuma epäonnistui: ${errorMsg}`;
        statusEl.style.color = "#FF8270";
      }
    } catch (err) {
      statusEl.textContent =
        "Tapahtuma epäonnistui: palvelimeen ei saatu yhteyttä.";
      statusEl.style.color = "#FF8270";
    }
  }
}

customElements.define("transaction-form", TransactionForm);
