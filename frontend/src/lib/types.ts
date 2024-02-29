export const enum FileEvent {
  Start = "start", // file upload started
  Finish = "finish", // file upload finished
  Progress = "progress", // a upload progress event
  Hash = "hash", // hashing started
  Error = "error", // any error
  Abort = "abort", // user aborted file
  MetaCreate = "meta_create", // file is ready to start uploading
  MetaFound = "meta_found", // file exists in some form on the server
  MetaNotFound = "meta_notfound" // file does not exist on the server
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
