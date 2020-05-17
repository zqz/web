import Config from '../Config.js';
import Meta from './Meta.js';
import Hash from './Hash.js';

const Upload = (file) => {
  let xhr = new XMLHttpRequest();
  let callbacks = {};

  function on(callback, func) {
    callbacks[callback] = func;
  }

  function cb(name, ...args) {
    const cb = callbacks[name];
    if ( cb === undefined) {
      return;
    }

    cb(...args);
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

    xhr.open('POST', `${Config.url}/api/file/` + file.hash, true);

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
    cb('error');
  }

  function onUploadAbort() {
    cb('abort');
  }

  function onUploadFinish(meta) {
    cb('finish', meta);
  }

  function onUploadStart() {
    cb('start');
  }

  function onHash() {
    cb('hash');
  };

  function onUploadProgress(event) {
    let prg = {
      loaded: getOffset() + event.loaded,
      total: file.data.size,
      time: now()
    };

    cb('progress', prg);
  }

  function onMetaFound(m) {
    file.meta = m;
    cb('meta_found', m)
  }

  function onMetaNotFound() {
    cb('meta_notfound');
  }

  function onMetaCheck() {
    cb('meta_check');
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
