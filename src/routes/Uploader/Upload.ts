import URLs from '$lib/urls';
import CallbacksHandler from './CallbacksHandler.js';
import Meta from './Meta.js';
import hashFile from '$lib/hash';

const UploadCallbacks = (file) => {
  let callbacks = CallbacksHandler();

  function on(...args) { callbacks.on(...args) }
  function onUploadError() { callbacks.call('error') }
  function onUploadAbort() { callbacks.call('abort') }
  function onUploadFinish(meta) { callbacks.call('finish', meta) }
  function onUploadStart() { callbacks.call('start') }
  function onUploadProgress(e) { callbacks.call('progress', e) }
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
    on
  }
}

const Upload = (file) => {
  let xhr = new XMLHttpRequest();
  let callbacks = UploadCallbacks(file);

  function abort() {
    xhr.abort();
  }

  function on(...args) {
    callbacks.on(...args);
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
    xhr.upload.addEventListener('progress', onUploadProgress);
    xhr.upload.addEventListener('error', callbacks.onUploadError);
    xhr.upload.addEventListener('abort', callbacks.onUploadAbort);
    xhr.addEventListener('readystatechange', onUploadStateChange);
    xhr.open('POST', URLs.postFileUrl(file.hash), true);
    xhr.send(file.data.slice(getOffset()));
    callbacks.onUploadStart();
  }

  function onUploadProgress(event) {
    let prg = {
      loaded: getOffset() + event.loaded,
      total: file.data.size,
      time: now()
    };

    callbacks.onUploadProgress(prg);
  }

  function onUploadStateChange(event) {
    if (xhr.readyState !== XMLHttpRequest.DONE) {
      return;
    }

    let text = event.target.responseText;
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
    on,
    fetchMeta,
    start,
    hash,
    abort,
  }
}

export default Upload;
