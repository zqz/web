export const enum FileEvent {
  Start = "start",
  Finish = "finish",
  Progress = "progress",
  Hash = "hash",
  Error = "error",
  Abort = "abort",
  MetaCreate = "meta_create",
  MetaCheck = "meta_check",
  MetaFound = "meta_found",
  MetaNotFound = "meta_notfound"
}

export interface FileProgress {
  loaded: number;
  total: number;
  time: number; // why is this a number?
}
