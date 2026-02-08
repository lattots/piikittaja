import { fetchUserTransactions } from "./api";
import { ApiError } from "./errors";
import { Transaction } from "./models";
import { format } from "./monetaryUtil";

export class TransactionTable extends HTMLElement {
	constructor() {
		super()
	}

	async connectedCallback() {
		this.innerHTML = `<p>Ladataan maksutapahtumia...</p>`
		try {
			await this.render();
		} catch (err) {
			this.innerHTML = `<p>Maksutapahtumien lataaminen epäonnistui.</p>`;
			console.error(err);
		}
	}

	async render() {
		const userId = this.getAttribute("user-id");
		if (!userId) {
			console.error("user-id not provided for transaction-table")
			return
		}
		const apiUrl = this.getAttribute("api-url");
		if (!apiUrl) {
			console.error("api-url not provided for transaction-table")
			return
		}

		const transactionQuantity: number = 30;

		let transactions: Transaction[];
		try {
			transactions = await fetchUserTransactions(apiUrl, userId, transactionQuantity);
		} catch (error) {
			if (error instanceof ApiError) {
				window.location.href = '/login';
				return;
			}
			console.log(error);
			return;
		}

		if (transactions.length === 0) {
			this.innerHTML = "<p>Tällä käyttäjällä ei ole vielä yhtäkään maksutapahtumaa.</p>"
			return;
		}

		this.innerHTML = `
            <table id="transaction-table">
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

	const now = new Date().getTime();
	const transactionTime = transaction.issuedAt.getTime();
	const isNew = (now - transactionTime) < 3000; // 3 s

	let flashClass = "";
	if (isNew) {
		flashClass = (transaction.type === "deposit") ? "new-deposit" : "new-withdraw";
	}

	return `
        <tr class="${flashClass}">
            <td>${formatDate(transaction.issuedAt)}</td>
            <td style="text-align: right; padding-right: 0.5rem">${format(transaction.amount)}</td>
            <td style="display: flex; justify-content: center; align-content: center; min-width: 32px">
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
