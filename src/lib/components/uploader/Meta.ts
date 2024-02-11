import URLs from '$lib/urls';
import { FileEvent, type Meta, type FileMetaRequest, type Uploadable } from '$lib/types';
import CallbacksHandler from './CallbacksHandler.js';

const MetaCallbacks = () => {
  let callbacks = CallbacksHandler<FileEvent>();

  function onMetaFound(meta: Meta) { callbacks.call(FileEvent.MetaFound, meta) }
  function onMetaNotFound() { callbacks.call(FileEvent.MetaNotFound) }
  function onMetaCreate(meta: Meta) { callbacks.call(FileEvent.MetaCreate, meta) }

  return {
    onMetaFound,
    onMetaNotFound,
    onMetaCreate,
    on: callbacks.on
  }
}

function buildFileMetaRequest(f: Uploadable) : FileMetaRequest {
  const d = f.data;

  if (f.hash === undefined) {
    throw "file missing hash";
  }

  return {
    name: d.name,
    type: d.type,
    path: '?',
    size: d.size,
    hash: f.hash
  }
}

const fetchFileMeta = (file: Uploadable) => {
  let meta:Meta;
  let callbacks = MetaCallbacks();
  let xhr = new XMLHttpRequest();

  function payload() {
    const r = buildFileMetaRequest(file);
    return JSON.stringify(r);
  }

  // actions
  async function retrieve() {
    const response = await fetch(URLs.getMetaUrl(file.hash!));
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

  function onMetaStateChange(e: Event) {
    if (xhr.readyState !== XMLHttpRequest.DONE) {
      return;
    }

    const target = e.target as XMLHttpRequest;
    const text = target.responseText;
    if (text === undefined || text === '') {
      return;
    }

    const meta = JSON.parse(text);
    callbacks.onMetaCreate(meta);
  }

  return {
    retrieve,
    create,
    on: callbacks.on,
  }
}

export default fetchFileMeta;
