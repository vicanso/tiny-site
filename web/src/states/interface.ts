export interface IList<T> {
  processing: boolean;
  items: T[];
  count: number;
}

export interface IStatus {
  desc: string;
  value: number;
}
