export class TransactionForm extends HTMLElement {
	constructor() {
		super();
	}

	connectedCallback() {
		this.innerHTML = `<p>Loading transaction form...</p>`
		this.render();
	}

	static get observedAttributes() {
		return ['api-url', 'user-id'];
	}

	render() {
		this.innerHTML = `
			<div>
				<form id="tx-form">
					<div>
						<label>Summa</label>
						<input type="number" name="amount" step="0.01" required placeholder="esim. 6.70">
					</div>
					<div>
						<button type="submit" name="type" value="deposit" class="submit-button">Talleta</button>
						<button type="submit" name="type" value="withdraw" class="submit-button">Nosta</button>
					</div>
				</form>
				<div id="status"></div>
			</div>
		`;

		const txFormEl = this.querySelector("#tx-form")
		if (!txFormEl) {
			console.log("No tx-form element")
			this.innerHTML = `<p>Something went wrong:(`
			return
		}
		txFormEl.addEventListener("submit", (e: Event) => this.handleSubmit(e));
	}

	async handleSubmit(event: Event) {
		if (!event) return
		event.preventDefault()

		const submitEvent = event as SubmitEvent
		const submitter = submitEvent.submitter as HTMLButtonElement
		const type = submitter.value; // "deposit" or "withdraw"

		const form = event.target as HTMLFormElement
		if (!form) return

		const statusEl = this.querySelector("#status") as HTMLElement
		if (!statusEl) return

		const formData = new FormData(form);
		const amountStr = formData.get("amount") as string
		if (!amountStr) {
			console.log("No amount in form")
			statusEl.innerHTML = "<p>Something went wrong:("
			return
		}
		const amount = parseFloat(amountStr);
		const amountInCents = Math.round(amount * 100)

		statusEl.textContent = "Sending..."
		statusEl.style.color = "#666"

		const apiUrl = this.getAttribute("api-url")
		const userId = this.getAttribute("user-id")

		try {
			const response = await fetch(`${apiUrl}/users/${userId}/transactions`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					amount: amountInCents,
					type: type,
				})
			});

			const errorMsg: string = await response.text()
			console.log(errorMsg)

			if (response.ok) {
				statusEl.textContent = `Tapahtuma onnistui!`
				statusEl.style.color = "#54DF60"
				form.reset();
			} else if (response.status === 402) {
				statusEl.textContent = "Rahat on loppu:("
				statusEl.style.color = "#FF8270"
			} else {
				statusEl.textContent = `Tapahtuma epäonnistui: ${errorMsg}`
				statusEl.style.color = "#FF8270"
			}
		} catch (err) {
			statusEl.textContent = "Tapahtuma epäonnistui: palvelimeen ei saatu yhteyttä."
			statusEl.style.color = "#FF8270"
		}
	}
}

customElements.define('transaction-form', TransactionForm)
