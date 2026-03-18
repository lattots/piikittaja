import { fetchUsers } from "./api.ts";

import tableStyles from "../css/table.css";
import generalStyles from "../css/general.css";

const tableSheet = new CSSStyleSheet();
tableSheet.replaceSync(tableStyles);

const generalSheet = new CSSStyleSheet();
generalSheet.replaceSync(generalStyles);

export class SearchBar extends HTMLElement {
  private timer: number | undefined;
  private abortController: AbortController | null = null;
  private inputElement: HTMLInputElement | null = null;

  private shadow: ShadowRoot;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });

    this.shadow.adoptedStyleSheets = [generalSheet, tableSheet];
  }

  connectedCallback() {
    this.render();
    this.setupEventListeners();
  }

  disconnectedCallback() {
    if (this.timer) clearTimeout(this.timer);
  }

  render() {
    this.shadow.innerHTML = `
			<div class="search-container">
				<input type="text" placeholder="Hae" />
			</div>
		`;
    this.inputElement = this.shadow.querySelector("input");
  }

  setupEventListeners() {
    if (!this.inputElement) return;

    this.inputElement.addEventListener("input", (e: Event) => {
      const target = e.target as HTMLInputElement;
      this.onSearchInput(target.value);
    });
  }

  onSearchInput(searchTerm: string) {
    clearTimeout(this.timer);

    if (this.abortController) {
      this.abortController.abort();
    }

    this.timer = window.setTimeout(() => {
      this.executeSearch(searchTerm);
    }, 500);
  }

  async executeSearch(searchTerm: string) {
    const apiUrl = this.getAttribute("api-url") || "";

    this.abortController = new AbortController();
    const { signal } = this.abortController;

    try {
      const users = await fetchUsers(apiUrl, searchTerm, signal);

      this.dispatchEvent(
        new CustomEvent("search-success", {
          detail: { users },
          bubbles: true,
          composed: true,
        }),
      );
    } catch (error: any) {
      if (error.name === "AbortError") return;

      this.dispatchEvent(
        new CustomEvent("search-error", {
          detail: { error },
          bubbles: true,
          composed: true,
        }),
      );
      console.log("Failed to search users");
    }
  }
}

customElements.define("search-bar", SearchBar);
