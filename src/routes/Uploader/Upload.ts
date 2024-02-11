import URLs from '$lib/urls';
import CallbacksHandler from './CallbacksHandler.js';
import Meta from './Meta.js';
import hashFile from '$lib/hash';
import type { FileEvent, FileProgress } from './types.js';

interface UploadHandler {
  onUploadError() : void;
  onUploadAbort() : void;
  onUploadProgress(e: FileProgress) : void;
  onUploadStart() : void;
  onUploadFinish() : void;
  onHash() : void;
  onMetaFound() : void;
  onMetaNotFound() : void;
  onMetaCheck() : void;
  on: (name: FileEvent, fn: Function) => void;
}

function UploadCallbacks(file) : UploadHandler {
  let callbacks = CallbacksHandler<FileEvent>();

  function onUploadError() { callbacks.call('error') }
  function onUploadAbort() { callbacks.call('abort') }
  function onUploadFinish(meta) { callbacks.call('finish', meta) }
  function onUploadStart() { callbacks.call('start') }
  function onUploadProgress(e: ProgressEvent) { callbacks.call('progress', e) }
  function onHash() { callbacks.call('hash') }

  function onMetaFound(m) {
    file.meta = m;
    callbacks.call('meta_found', m)
  }

  function onMetaNotFound() { callbacks.call('meta_notfound') }
  function onMetaCheck() { callbacks.call('meta_check') }

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

export const uploadFile = (file) => {
  let xhr = new XMLHttpRequest();
  let callbacks = UploadCallbacks(file);

  function abort() {
    xhr.abort();
  }

  function start() {
    if (file.meta !== undefined) {
      upload();
      return;
    }

    let m = Meta(file);
    m.on('create', function() {
      file.meta = m.get();
      upload();
    });
    m.create();
  }

  function fetchMeta() {
    let m = Meta(file);
    m.on('found', callbacks.onMetaFound);
    m.on('notfound', callbacks.onMetaNotFound);
    m.retrieve();
  }

  async function hash() {
    callbacks.onHash();
    hashFile(file.data, function(h) {
      file.hash = h;
      fetchMeta();
    });
  }

  function upload() {
    xhr.upload.addEventListener('progress', onProgress);
    xhr.upload.addEventListener('error', callbacks.onUploadError);
    xhr.upload.addEventListener('abort', callbacks.onUploadAbort);
    xhr.addEventListener('readystatechange', onStateChange);
    xhr.open('POST', URLs.postFileUrl(file.hash), true);
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
    callbacks.onUploadFinish(meta);
  }

  function getOffset() {
    const m = file.meta;

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

