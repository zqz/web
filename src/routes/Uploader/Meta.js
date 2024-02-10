import Config from '../Config';

import CallbacksHandler from './CallbacksHandler.js';

const MetaCallbacks = () => {
  let callbacks = CallbacksHandler();

  function on(...args) { callbacks.on(...args) }
  function onMetaFound(meta) { callbacks.call('found', meta) }
  function onMetaNotFound() { callbacks.call('notfound') }
  function onMetaCreate() { callbacks.call('create') }

  return {
    onMetaFound,
    onMetaNotFound,
    onMetaCreate,
    on
  }
}

const Meta = (file) => {
  let meta = null;
  let callbacks = MetaCallbacks();
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

  function on(...args) {
    callbacks.on(...args);
  }

  async function retrieve() {
    const response = await fetch(Config.getMetaUrl(file.hash));
    const json = await response.json();

    if (json.message === 'file not found') {
      callbacks.onMetaNotFound();
      return;
    }

    callbacks.onMetaFound(json);
  }

  function create() {
    xhr.addEventListener('readystatechange', onMetaStateChange);
    xhr.open('POST', Config.postMetaUrl(), true);
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
    callbacks.onMetaCreate();
  }

  return {
    retrieve,
    create,
    on,
    get
  }
}

export default Meta;
