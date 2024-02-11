import URLs from '$lib/urls';
import CallbacksHandler from './CallbacksHandler.js';
import hashFile from '$lib/hash';
import { FileEvent } from './types.js';
import type { FileProgress, Meta, Uploadable } from './types.js';
import fetchFileMeta from './Meta.js';

interface UploadHandler {
  onUploadError() : void;
  onUploadAbort() : void;
  onUploadProgress(e: FileProgress) : void;
  onUploadStart() : void;
  onUploadFinish() : void;
  onHash() : void; // start hashing
  onMetaFound() : void;
  onMetaNotFound() : void;
  onMetaCheck() : void;
  on: (name: FileEvent, fn: Function) => void;
}

function UploadCallbacks(file: Uploadable) : UploadHandler {
  let callbacks = CallbacksHandler<FileEvent>();

  function onUploadError() { callbacks.call(FileEvent.Error) }
  function onUploadAbort() { callbacks.call(FileEvent.Abort) }
  function onUploadFinish() { callbacks.call(FileEvent.Finish) }
  function onUploadStart() { callbacks.call(FileEvent.Start) }
  function onUploadProgress(e: FileProgress) { callbacks.call(FileEvent.Progress, e) }
  function onHash() { callbacks.call(FileEvent.Hash) }

  function onMetaFound(m: Meta) {
    console.log('found meta');
    file.meta = m;
    console.log('assigned meta', m);
    callbacks.call(FileEvent.MetaFound);
  }

  function onMetaNotFound() { callbacks.call(FileEvent.MetaNotFound) }
  function onMetaCheck() { callbacks.call(FileEvent.MetaCheck) }

  return {
    onUploadAbort,
    onUploadProgress,
    onUploadError,
    onUploadStart,
    onUploadFinish,
    onHash,
    onMetaCheck,
    onMetaFound,
    onMetaNotFound,
    on: callbacks.on,
  }
}

export const uploadFile = (file: Uploadable) => {
  let xhr = new XMLHttpRequest();
  let callbacks = UploadCallbacks(file);
  
    let m = fetchFileMeta(file);
    m.on(FileEvent.MetaFound, callbacks.onMetaFound);
    m.on(FileEvent.MetaNotFound, callbacks.onMetaNotFound);
    m.on(FileEvent.MetaCreate, (meta: Meta) => {
      file.meta = meta;
      console.log('starting: meta created, uploading', meta);
      upload()
    });

  function abort() {
    xhr.abort();
  }

  function start() {
    console.log('starting');
    // if there is a metadata, we can can start uploading
    if (file.meta !== undefined) {
      console.log('starting: meta existing');
      upload();
      return;
    }

    console.log('starting: creating meta');
    m.create();
  }

  function fetchMeta() {
    console.log('fetching meta');
    m.retrieve();
  }

  async function hash() {
    console.log('hashing file');
    callbacks.onHash();

    hashFile(file.data, function(h) {
      console.log('finished hashing file');
      file.hash = h;
      fetchMeta();
    });
  }

  function upload() {
    xhr.upload.addEventListener('progress', onProgress);
    xhr.upload.addEventListener('error', callbacks.onUploadError);
    xhr.upload.addEventListener('abort', callbacks.onUploadAbort);
    xhr.addEventListener('readystatechange', onStateChange);
    xhr.open('POST', URLs.postFileUrl(file.hash!), true);
    xhr.send(file.data.slice(getOffset()));
    callbacks.onUploadStart();
  }

  function onProgress(event: ProgressEvent) {
    let fileProgress = {
      loaded: getOffset() + event.loaded,
      total: file.data.size,
      time: now()
    };

    callbacks.onUploadProgress(fileProgress);
  }

  function onStateChange(event: Event) {
    if (xhr.readyState !== XMLHttpRequest.DONE) {
      return;
    }

    const target = event.target as XMLHttpRequest;
    let text = target.responseText;
    if (text === undefined || text === '') {
      return;
    }

    const meta = JSON.parse(text);
    file.meta = meta;
    callbacks.onUploadFinish();
  }

  function getOffset() {
    const m = file.meta!;

    if (m.bytes_received === undefined) {
      return 0;
    }

    return m.bytes_received;
  }

  function now() {
    const d = new Date();
    return d.getTime();
  }

  return {
    on: callbacks.on,
    fetchMeta,
    start,
    hash,
    abort,
  }
}

