export const enum FileEvent {
  Start = "start",
  Finish = "finish",
  Progress = "progress",
  Hash = "hash",
  MetaCheck = "meta_check",
  MetaFound = "meta_found",
  MetaNotFound = "meta_notfound"
}

export type FileProgress = {
  loaded: number;
  total: number;
  time: number; // why is this a number?
}
