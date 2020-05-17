import Config from '../Config.js';

const Meta = (file) => {
  let meta = null;
  let callbacks = {};
  let xhr = new XMLHttpRequest();

  function payload() {
    const d = file.data;

    let p = {
      name: d.name,
      type: d.type,
      path: d.path,
      size: d.size,
      hash: file.hash
    };

    return JSON.stringify(p);
  }

  // actions
  function get() {
    return meta;
  }

  async function retrieve() {
    const response = await fetch(`${Config.url}/api/meta/by-hash/${file.hash}`);
    const json = await response.json();

    if (json.message === 'file not found') {
      onMetaNotFound();
      return;
    }

    onMetaFound(json);
  }

  function create() {
    xhr.addEventListener('readystatechange', onMetaStateChange);
    xhr.open('POST', `${Config.url}/api/meta`, true);
    xhr.send(payload());
  }

  function onMetaStateChange(e) {
    if (xhr.readyState !== XMLHttpRequest.DONE) {
      return;
    }

    let text = e.target.responseText;
    if (text === undefined || text === '') {
      return;
    }

    meta = JSON.parse(text);
    onMetaCreate();
  }

  // callbacks
  function onMetaFound(meta) {
    cb('found', meta);
  }

  function onMetaNotFound() {
    cb('notfound');
  }

  function onMetaCreate() {
    cb('create');
  }

  // callback helpers
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

  return {
    retrieve,
    create,
    on,
    get
  }
}

export default Meta;
