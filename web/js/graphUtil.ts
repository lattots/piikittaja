import { Transaction } from "./models.ts";

export interface ChartData {
  dates: string[];
  values: number[];
  maxVal: number;
}

export function processTransactionsForGraph(
  transactions: Transaction[],
  endDate: Date,
  windowDays: number,
): ChartData {
  const sumsByDate = new Map<string, number>();

  transactions.forEach((t) => {
    if (t.type !== "withdraw") return;

    const dateKey = t.issuedAt.toISOString().split("T")[0];
    const current = sumsByDate.get(dateKey) || 0;
    sumsByDate.set(dateKey, current + t.amount);
  });

  const dates: string[] = [];
  const values: number[] = [];
  let maxVal = 0;

  for (let i = windowDays - 1; i >= 0; i--) {
    const d = new Date(endDate);
    d.setDate(d.getDate() - i);

    const dateKey = d.toISOString().split("T")[0];
    const val = sumsByDate.get(dateKey) || 0; // Default to 0 if no data

    dates.push(dateKey);
    values.push(val);

    if (val > maxVal) maxVal = val;
  }

  return { dates, values, maxVal };
}
