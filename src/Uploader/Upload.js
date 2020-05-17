import Config from '../Config.js';
import Meta from './Meta.js';
import Hash from './Hash.js';

const CallbackHandler = () => {
  let callbacks = {};

  function on(callback, func) {
    callbacks[callback] = func;
  }

  function call(name, ...args) {
    const cb = callbacks[name];
    if ( cb === undefined) {
      return;
    }

    cb(...args);
  }

  return { on, call };
}

const Upload = (file) => {
  let xhr = new XMLHttpRequest();
  let callbacks = CallbackHandler();

  function on(...args) {
    callbacks.on(...args);
  }

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
    m.on('found', onMetaFound);
    m.on('notfound', onMetaNotFound);
    m.retrieve();
  }

  async function hash() {
    onHash();
    Hash(file, function(h) {
      file.hash = h;
      fetchMeta();
    });
  }

  function upload() {
    xhr.upload.addEventListener('progress', onUploadProgress);
    xhr.upload.addEventListener('error', onUploadError);
    xhr.upload.addEventListener('abort', onUploadAbort);
    xhr.addEventListener('readystatechange', onUploadStateChange);
    xhr.open('POST', Config.postFileUrl(file.hash), true);
    xhr.send(file.data.slice(getOffset()));
    onUploadStart();
  }

  function getOffset() {
    const m = file.meta;

    if (m.bytes_received === undefined) {
      return 0;
    }

    return m.bytes_received;
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
    onUploadFinish(meta);
  }

  function onUploadError() {
    callbacks.call('error');
  }

  function onUploadAbort() {
    callbacks.call('abort');
  }

  function onUploadFinish(meta) {
    callbacks.call('finish', meta);
  }

  function onUploadStart() {
    callbacks.call('start');
  }

  function onHash() {
    callbacks.call('hash');
  };

  function onUploadProgress(event) {
    let prg = {
      loaded: getOffset() + event.loaded,
      total: file.data.size,
      time: now()
    };

    callbacks.call('progress', prg);
  }

  function onMetaFound(m) {
    file.meta = m;
    callbacks.call('meta_found', m)
  }

  function onMetaNotFound() {
    callbacks.call('meta_notfound');
  }

  function onMetaCheck() {
    callbacks.call('meta_check');
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
