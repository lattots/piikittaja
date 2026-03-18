import statsStyles from "../css/stats.css";
import generalStyles from "../css/general.css";

interface WindowChangedDetail {
  windowDays: number;
}

const statsSheet = new CSSStyleSheet();
statsSheet.replaceSync(statsStyles);

const generalSheet = new CSSStyleSheet();
generalSheet.replaceSync(generalStyles);

export class StatControls extends HTMLElement {
  private shadow: ShadowRoot;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });

    this.shadow.adoptedStyleSheets = [generalSheet, statsSheet];
  }

  connectedCallback(): void {
    this.render();
    this.setupEventListeners();
  }

  private render(): void {
    this.shadow.innerHTML = `
      <form id="stat-control-form">
        <label class="window-control-label button"> 30 pv <input type="radio" name="window" value="30" checked /></label>
        <label class="window-control-label button"> 3 kk <input type="radio" name="window" value="90" /></label>
        <label class="window-control-label button"> 6 kk <input type="radio" name="window" value="180" /></label>
        <label class="window-control-label button"> 1 v <input type="radio" name="window" value="365" /></label>
      </form>
    `;
  }

  private setupEventListeners(): void {
    const form = this.shadow.querySelector<HTMLFormElement>(
      "#stat-control-form",
    );

    form?.addEventListener("change", (event: Event) => {
      const target = event.target as HTMLInputElement;
      const newWindow = parseInt(target.value, 10);

      const changeEvent = new CustomEvent<WindowChangedDetail>(
        "window-changed",
        {
          bubbles: true,
          composed: true,
          detail: { windowDays: newWindow },
        },
      );

      this.dispatchEvent(changeEvent);
    });
  }
}

customElements.define("stat-controls", StatControls);
