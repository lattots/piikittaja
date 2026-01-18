import { Transaction } from "./models";
import { format } from "./monetaryUtil";

export class TransactionTable extends HTMLElement {
	constructor() {
		super()
	}

	async connectedCallback() {
		this.innerHTML = `<p style="font-weight: bold">Loading transactions...</p>`
		await this.render()
	}

	async render() {
		const userId = this.getAttribute("user-id");
		const apiUrl = this.getAttribute("api-url");

		const transactionQuantity: number = 10;

		const resp = await fetch(`${apiUrl}/users/${userId}/transactions?quantity=${transactionQuantity}`)
		if (!resp.ok) return;

		const rawData: any[] = await resp.json();

		if (!rawData || rawData.length === 0) {
			this.innerHTML = "<p>Tällä käyttäjällä ei ole vielä yhtäkään maksutapahtumaa</p>"
			return
		}

		const transactions: Transaction[] = rawData.map(t => ({
			...t,
			issuedAt: new Date(t.issuedAt)
		}));

		this.innerHTML = `
            <table>
                <thead style="font-weight: bold">
                    <tr>
						<th>Aika</th>
						<th>Summa</th>
					</tr>
                </thead>
                <tbody>
                    ${transactions.map(renderTransaction).join("")}
                </tbody>
            </table>
        `;
	}
}

function renderTransaction(transaction: Transaction): string {
	const assetDir: string = "/app/assets/";
	const typeIconSource: string = (transaction.type === "deposit") ? "deposit.svg" : "withdraw.svg";
	const iconColor: string = (transaction.type === "deposit") ? "#54DF60" : "#FF8270";
	const iconPath: string = assetDir + typeIconSource;

	return `
        <tr>
            <td>${formatDate(transaction.issuedAt)}</td>
            <td style="text-align: right; padding-right: 8px">${format(transaction.amount)}</td>
            <td style="max-width: 32px; min-width: 32px">
                <div style="
                    width: 24px; 
                    height: 24px; 
                    background-color: ${iconColor}; 
                    -webkit-mask: url('${iconPath}') no-repeat center;
                    mask: url('${iconPath}') no-repeat center;
                    mask-size: contain;
                    -webkit-mask-size: contain;
                "></div>
            </td>
        </tr>
    `;
}

function formatDate(date: Date): string {
	const day = String(date.getDate()).padStart(2, "0");
	const month = String(date.getMonth() + 1).padStart(2, "0"); // Months are 0-indexed
	const year = String(date.getFullYear()).slice(-2); // Get last 2 digits
	const hours = String(date.getHours()).padStart(2, "0");
	const minutes = String(date.getMinutes()).padStart(2, "0");

	return `${day}.${month}.${year} ${hours}:${minutes}`;
}

customElements.define("transaction-table", TransactionTable);
