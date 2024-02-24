import { URLs, hashFile } from '$lib/util';
import { callbacks } from './callbacks';
import { FileEvent } from '$lib/types';
import type { Meta, Uploadable } from '$lib/types';
import { fetchFileMeta } from './fetchMeta';

export const uploadFile = (file: Uploadable) => {
  let xhr = new XMLHttpRequest();
  let cb = callbacks<FileEvent>();
  let m = fetchFileMeta(file);

  m.on(FileEvent.MetaFound, (m: Meta) => {
    file.meta = m;
    cb.emit(FileEvent.MetaFound);
  });

  m.on(FileEvent.MetaNotFound, () => cb.emit(FileEvent.MetaNotFound))
  m.on(FileEvent.MetaCreate, (meta: Meta) => {
    file.meta = meta;
    upload()
  });

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
    xhr.send(file.data.slice(getOffset()));
    cb.emit(FileEvent.Start);
  }

  function onProgress(event: ProgressEvent) {
    let fileProgress = {
      loaded: getOffset() + event.loaded,
      total: file.data.size,
      time: now()
    };

    cb.emit(FileEvent.Progress, fileProgress);
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

  function getOffset() {
    const m = file.meta!;

    if (m.bytes_received === undefined) {
      return 0;
    }

    return m.bytes_received;
  }

  return {
    on: cb.on,
    fetchMeta,
    start,
    hash,
    abort,
  }
}

function now() {
  const d = new Date();
  return d.getTime();
}

