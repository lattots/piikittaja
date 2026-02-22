import { fetchTransactions } from "./api.ts";

import {
  CategoryScale,
  Chart,
  Filler,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  Title,
  Tooltip,
} from "chart.js";
import { processTransactionsForGraph } from "./graphUtil.ts";

Chart.register(
  LineController,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Title,
  Tooltip,
  Filler,
);

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
    this.canvas = this.shadowRoot!.querySelector(
      "#chart-canvas",
    ) as HTMLCanvasElement;
  }

  async connectedCallback() {
    await this.render();
  }

  private async render() {
    const apiUrl = this.getAttribute("api-url") || "";
    const endDate = new Date();

    const rawData = await fetchTransactions(
      apiUrl,
      endDate,
      this.windowDays,
      "withdraw",
    );
    const { dates, values } = processTransactionsForGraph(
      rawData,
      endDate,
      this.windowDays,
    );

    if (this.chart) this.chart.destroy();

    this.chart = new Chart(this.canvas, {
      type: "line",
      data: {
        labels: dates,
        datasets: [{
          data: values,
          borderColor: "rgb(255, 130, 112)",
          backgroundColor: "rgba(255, 130, 112, 0.15)",
          fill: true,
          tension: 0.3,
          pointRadius: this.windowDays > 60 ? 0 : 3,
        }],
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
              color: "rgba(255, 255, 255, 1)",
              callback: (val: any, index: any) => {
                const d = new Date(dates[index]);
                return d.toLocaleString("fi-FI", {
                  month: "short",
                  day: "numeric",
                });
              },
            },
          },
          y: {
            beginAtZero: true,
            ticks: {
              color: "rgba(255, 255, 255, 1)",
              callback: (value: any) => `${(Number(value) / 100).toFixed(0)} €`,
            },
          },
        },
        plugins: {
          tooltip: {
            callbacks: {
              label: (ctx: any) =>
                `Summa: ${(ctx.parsed.y / 100).toLocaleString()} €`,
            },
          },
        },
      },
    });
  }
}

customElements.define("app-stats", AppStats);
