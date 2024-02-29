import type { FileProgress, Meta, Uploadable } from '$lib/types';
import { FileEvent } from '$lib/types';
import { URLs, hashFile } from '$lib/util';
import { callbacks } from './callbacks';
import { fetchFileMeta } from './fetchMeta';

export const uploadFile = (file: Uploadable) => {
  let xhr = new XMLHttpRequest();
  let cb = callbacks<FileEvent>();
  let m = fetchFileMeta(file);

  m.on(FileEvent.MetaFound, (meta: Meta) => {
    setFileMetadata(meta);
    cb.emit(FileEvent.MetaFound);
  });

  m.on(FileEvent.MetaNotFound, () => cb.emit(FileEvent.MetaNotFound))
  m.on(FileEvent.MetaCreate, (meta: Meta) => {
    setFileMetadata(meta);
    upload()
  });

  function setFileMetadata(m: Meta) {
    file.meta = m;
  }

  function abort() {
    xhr.abort();
  }

  function start() {
    if (file.meta) {
      upload();
    } else {
      m.create();
    }
  }

  function fetchMeta() {
    m.retrieve();
  }

  async function hash() {
    cb.emit(FileEvent.Hash);

    try {
      hashFile(file.data, (h: string) => {
        file.hash = h;
        fetchMeta();
      });
    } catch(e) {
      console.log('caught error');
      cb.emit(FileEvent.Error);
    }
  }

  function upload() {
    xhr.upload.addEventListener('progress', onProgress);
    xhr.upload.addEventListener('error', () => cb.emit(FileEvent.Error));
    xhr.upload.addEventListener('abort', () => cb.emit(FileEvent.Abort));

    xhr.addEventListener('readystatechange', onStateChange);
    xhr.open('POST', URLs.postFileUrl(file.hash!), true);
    xhr.send(payloadData());
    cb.emit(FileEvent.Start);
  }

  function payloadData() : Blob {
    const offset = getMetaOffset(file.meta!);
    return file.data.slice(offset);
  }

  function onProgress(event: ProgressEvent) {
    cb.emit(FileEvent.Progress, buildProgressEvent(file.meta!, event));
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
    file.meta = meta;
    cb.emit(FileEvent.Finish)
  }

  return {
    on: cb.on,
    fetchMeta,
    start,
    hash,
    abort,
  }
}

function getMetaOffset(m: Meta) : number {
  if (m === undefined || m.bytes_received === undefined) {
    return 0;
  }

  return m.bytes_received;
}

function now() {
  const d = new Date();
  return d.getTime();
}

function buildProgressEvent(m: Meta, e: ProgressEvent) : FileProgress {
  const offset = getMetaOffset(m);

  return {
    loaded: offset + e.loaded,
    total: m.size,
    time: now()
  };
}

