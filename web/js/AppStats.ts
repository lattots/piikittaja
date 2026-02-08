import { fetchTransactions } from "./api";

import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, Tooltip, Filler } from "chart.js";
import { processTransactionsForGraph } from "./graphUtil";

Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Title, Tooltip, Filler);

export class AppStats extends HTMLElement {
	private chart: Chart | null = null;
	private canvas: HTMLCanvasElement;
	windowDays: number = 30;

	constructor() {
		super();
		this.attachShadow({ mode: "open" });
		this.shadowRoot!.innerHTML = `
			<canvas id="chart-canvas"></canvas>
        `;
		this.canvas = this.shadowRoot!.querySelector("#chart-canvas") as HTMLCanvasElement;
	}

	async connectedCallback() {
		await this.render();
	}

	private async render() {
		const apiUrl = this.getAttribute('api-url') || '';
		const endDate = new Date();

		// 1. Fetch & Process data (Uses the same logic as before to fill gaps with 0)
		const rawData = await fetchTransactions(apiUrl, endDate, this.windowDays, 'withdraw');
		const { dates, values } = processTransactionsForGraph(rawData, endDate, this.windowDays);

		// 2. Destroy old chart if it exists (for dynamic updates)
		if (this.chart) this.chart.destroy();

		// 3. Initialize Chart.js
		this.chart = new Chart(this.canvas, {
			type: 'line',
			data: {
				labels: dates,
				datasets: [{
					label: 'Withdrawals (Cents)',
					data: values,
					borderColor: '#4A90E2',
					backgroundColor: 'rgba(74, 144, 226, 0.1)',
					fill: true,
					tension: 0.3,
					pointRadius: this.windowDays > 60 ? 0 : 3,
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					x: {
						grid: { display: false },
						ticks: {
							autoSkip: true,
							maxRotation: 0,
							callback: (val, index) => {
								const d = new Date(dates[index]);
								return d.toLocaleString('default', { month: 'short', day: 'numeric' });
							}
						}
					},
					y: {
						beginAtZero: true,
						ticks: {
							callback: (value) => `${(Number(value) / 100).toFixed(0)} €`
						}
					}
				},
				plugins: {
					tooltip: {
						callbacks: {
							label: (ctx) => `Summa: ${(ctx.parsed.y / 100).toLocaleString()} €`
						}
					}
				}
			}
		});
	}
}

customElements.define("app-stats", AppStats);
