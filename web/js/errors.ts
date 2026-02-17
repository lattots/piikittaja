type ErrorName =
  | `UNAUTHENTICATED_ERROR`
  | `FORBIDDEN_ERROR`;

export class ApiError extends Error {
  override name: ErrorName;
  override message: string;

  constructor(name: ErrorName, message: string) {
    super();
    this.name = name;
    this.message = message;
  }
}
