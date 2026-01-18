export class AppStats extends HTMLElement {
	constructor() {
		super()
	}

	async connectedCallback() {
		this.innerHTML = `<p>Loading statistics...</p>`;
		try {
			await this.render();
		} catch (err) {
			this.innerHTML = `<p>Failed to load stats.</p>`;
			console.error(err);
		}
	}

	async render() {
		this.innerHTML = `<img src="/app/assets/otso.jpg" style="max-width: 200px; max-height: 200px; border-radius: 24px" alt="Nice Statistics:)">`
	}
}

customElements.define("app-stats", AppStats);
