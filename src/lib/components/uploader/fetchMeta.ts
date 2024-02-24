import { URLs } from '$lib/util';
import { FileEvent, type Meta, type FileMetaRequest, type Uploadable } from '$lib/types';
import { callbacks } from './callbacks';

export const fetchFileMeta = (file: Uploadable) => {
  let cb = callbacks<FileEvent>();
  let xhr = new XMLHttpRequest();

  // actions
  async function retrieve() {
    const response = await fetch(URLs.getMetaUrl(file.hash!), { mode: "cors" });
    if (response.status === 204) {
      cb.emit(FileEvent.MetaNotFound);
      return;
    }

    const json = await response.json();
    cb.emit(FileEvent.MetaFound, json);
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
    cb.emit(FileEvent.MetaCreate, meta);
  }

  function payload() {
    const r = buildFileMetaRequest(file);
    return JSON.stringify(r);
  }

  return {
    retrieve,
    create,
    on: cb.on,
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

