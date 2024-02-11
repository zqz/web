import URLs from '$lib/urls';
import { FileEvent } from './types';
import CallbacksHandler from './CallbacksHandler.js';

const MetaCallbacks = () => {
  let callbacks = CallbacksHandler<FileEvent>();

  function onMetaFound(meta) { callbacks.call(FileEvent.MetaFound, meta) }
  function onMetaNotFound() { callbacks.call(FileEvent.MetaNotFound) }
  function onMetaCreate() { callbacks.call(FileEvent.MetaCreate) }

  return {
    onMetaFound,
    onMetaNotFound,
    onMetaCreate,
    on: callbacks.on
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
    const response = await fetch(URLs.getMetaUrl(file.hash));
    const json = await response.json();

    if (json.message === 'file not found') {
      callbacks.onMetaNotFound();
      return;
    }

    callbacks.onMetaFound(json);
  }

  function create() {
    xhr.addEventListener('readystatechange', onMetaStateChange);
    xhr.open('POST', URLs.postMetaUrl(), true);
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
