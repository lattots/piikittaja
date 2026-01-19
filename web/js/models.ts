export type User = {
	id: number
	username: string
	firstName: string
	lastName: string
	balance: number // Balance in cents
}

export type Transaction = {
	issuer: string
	issuedAt: Date
	type: string
	amount: number // Transaction amount in cents
}
