export function format(monetaryValue: number): string {
  const euros: number = monetaryValue / 100;
  return `${euros} €`;
}
