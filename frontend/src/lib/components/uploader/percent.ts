import type { FileProgress } from "$lib/types";

export function calcPercent(updates: Array<FileProgress>) : number {
  if (updates.length === 0) {
    return 0;
  }

  let lastUpdate = updates[updates.length-1];
  
  return lastUpdate.loaded / lastUpdate.total;
}
