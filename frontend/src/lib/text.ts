export function truncate(v: string, len: number) : string {
  if (v.length < len - 3) {
    return v;
  }

  return `${v.slice(0, len)}...`;
}

export function generateId() : string {
  return Math.random().toString(20).slice(2, 8);
}
