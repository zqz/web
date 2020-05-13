import Config from '../Config.js';

const Upload = () => {
  let file = null;
  let xhr = new XMLHttpRequest();
  let callbacks = {};

  function setFile(f) {
    file = f;
  }

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

  function onUploadProgress(event) {
    let prg = {
      loaded: getOffset() + event.loaded,
      total: file.data.size,
      time: now()
    };

    cb('progress', prg);
  }

  function now() {
    const d = new Date();
    return d.getTime();
  }

  return {
    on,
    upload,
    abort,
    setFile
  }
}

export default Upload;
