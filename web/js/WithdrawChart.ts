import { fetchTransactions } from "./api.ts";
import { ChartData } from "./graphUtil.ts";

import statsStyles from "../css/stats.css";
import generalStyles from "../css/general.css";

const statsSheet = new CSSStyleSheet();
statsSheet.replaceSync(statsStyles);

const generalSheet = new CSSStyleSheet();
generalSheet.replaceSync(generalStyles);

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

export class WithdrawChart extends HTMLElement {
  private chart: Chart | null = null;
  private canvas: HTMLCanvasElement;
  windowDays: number = 30;

  private shadow: ShadowRoot;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.innerHTML = `
			<canvas id="chart-canvas"></canvas>
        `;
    this.canvas = this.shadow.querySelector(
      "#chart-canvas",
    ) as HTMLCanvasElement;

    this.shadow.adoptedStyleSheets = [generalSheet, statsSheet];
  }

  async connectedCallback() {
    await this.render();
    this.setupEventListeners();
  }

  disconnectedCallback() {
    window.removeEventListener(
      "window-changed",
      this.handleWindowChanged as EventListener,
    );
  }

  async render() {
    const apiUrl = this.getAttribute("api-url") || "";
    const endDate = new Date();

    const rawData = await fetchTransactions(
      apiUrl,
      endDate,
      this.windowDays,
      "withdraw",
    );
    const chartData = processTransactionsForGraph(
      rawData,
      endDate,
      this.windowDays,
    );

    if (this.chart) this.chart.destroy();

    this.chart = createChart(chartData, this.canvas);
  }

  private setupEventListeners() {
    window.addEventListener("window-changed", this.handleWindowChanged);
  }

  private handleWindowChanged = async (e: Event) => {
    const customEvent = e as CustomEvent<{ windowDays: number }>;
    this.windowDays = customEvent.detail.windowDays;
    await this.render();
  };
}

function createChart(data: ChartData, canvas: HTMLCanvasElement): Chart {
  return new Chart(canvas, {
    type: "line",
    data: {
      labels: data.dates,
      datasets: [{
        data: data.values,
        borderColor: "rgb(255, 130, 112)",
        backgroundColor: "rgba(255, 130, 112, 0.15)",
        fill: true,
        tension: 0.3,
        pointRadius: data.windowDays > 60 ? 0 : 3,
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
              const d = new Date(data.dates[index]);
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

customElements.define("withdraw-chart", WithdrawChart);
