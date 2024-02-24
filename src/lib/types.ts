export const enum FileEvent {
  Start = "start",
  Finish = "finish",
  Progress = "progress",
  Hash = "hash",
  Error = "error",
  Abort = "abort",
  MetaCreate = "meta_create",
  MetaFound = "meta_found",
  MetaNotFound = "meta_notfound"
}

export const enum FileStatus {
  Queue = "queued",
  Hashing = "hashing",
  Ready = "ready",
  InProgress = "in_progress",
  Done = "done",
  Error = "error"
}

export interface FileProgress {
  loaded: number;
  total: number;
  time: number; // why is this a number?
}

export interface Uploadable {
  internalId: string;
  data: globalThis.File;
  meta?: Meta;
  hash?: string;
}

export interface Meta {
  bytes_received: number;
  date: string;
  hash: string;
  name: string;
  size: number;
  slug: string;
  type: string;

  alias: string;
  thumbnail: string;
}

export interface FileMetaRequest {
  name: string;
  type: string;
  path: string;
  size: number;
  hash: string;
}
