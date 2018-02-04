// var shajs = require('sha.js')
import Rusha from 'rusha';

class FileData {
  constructor(file) {
    this.file = file;
    this.state = null;
    this.onProgressCallback = null;
    this.onStartCallback = null;
    this.onAbortCallback = null;
    this.onErrorCallback = null;
    this.onFinishCallback = null;
    this.onResponseCallback = null;
    this.onHashCallback = null;
    this.onExistsCallback = null;
    this.xhr = null;
    this.token = null;
    this.bytesPreviouslySent = 0;
    this.bytesSent = 0;

    this._key = Math.random().toString(36).substring(7);
  }

  previouslySent = () => {
    return this.bytesPreviouslySent;
  }

  hash = () => {
    var worker = Rusha.createWorker();
    worker.onmessage = (e) => {
      this.file.hash = e.data.hash;
      this.check();
      this.onHashCallback(e.data);
    }

    worker.postMessage({id: this.key(), data: this.file});
  }

  key() {
    return this._key;
  }

  isResumable = () => {
    return this.bytesPreviouslySent > 0
  }

  onProgress(callback) {
    this.onProgressCallback = callback;
  }

  onError(callback) {
    this.onErrorCallback = callback;
  }

  onStart(callback) {
    this.onStartCallback = callback;
  }

  onAbort(callback) {
    this.onAbortCallback = callback;
  }

  onLoad(callback) {
    this.onLoadCallback = callback;
  }

  onExists(callback) {
    this.onExistsCallback = callback;
  }

  onResponse(callback) {
    this.onResponseCallback = callback;
  }

  onHash(callback) {
    this.onHashCallback = callback;
  }

  blob() {
    return this.file.blob;
  }

  started() {
    return this.state === 'started';
  }

  check() {
    console.log('checking existing state');
    var pxhr = new XMLHttpRequest();
    pxhr.addEventListener('readystatechange', (e) => {
      if (pxhr.readyState !== 4) {
        return;
      }

      if (pxhr.status !== 200) {
        return;
      }

      var text = e.target.responseText;
      if (text === undefined) {
        return;
      }

      var response = JSON.parse(text);

      this.bytesPreviouslySent = response.bytes_received;

      if (this.bytesPreviouslySent === this.file.size) {
        this.state = 'finished';
      }
      this.onExistsCallback(response);
    });

    pxhr.open('GET', 'http://localhost:3001/upload/' + this.file.hash, true);
    pxhr.send();
  }

  prepare() {
    var data = {
      name: this.file.name,
      hash: this.file.hash,
      size: this.file.size,
    }

    var pxhr = new XMLHttpRequest();
    pxhr.addEventListener('readystatechange', (e) => {
      if (pxhr.readyState === 4) {
        var text = e.target.responseText;

        if (text === undefined) {
          return;
        }

        var response = JSON.parse(text);
        console.log(response);
        this.token = response.token;
        this.bytesPreviouslySent = response.bytes_received;
        this.start();
      }
    });

    pxhr.open('POST', 'http://localhost:3001/prepare', true);
    pxhr.send(JSON.stringify(data));
  }

  start() {
    this.xhr = new XMLHttpRequest();
    this.xhr.upload.addEventListener('progress', (e) => {
      this.bytesSent = this.bytesPreviouslySent + e.loaded;
      this.onProgressCallback({
        loaded: this.bytesSent,
        total: this.file.size
      });
    });

    this.xhr.upload.addEventListener('load', () => {
      this.state = 'finished';
      this.onLoadCallback();
    });

    this.xhr.addEventListener('readystatechange', (e) => {
      if (this.xhr.readyState === 4) {
        var text = e.target.responseText;

        if (text === undefined || text === null || text.length === 0) {
          return;
        }

        var data = JSON.parse(text);
        this.onResponseCallback(data);
      }
    });

    this.xhr.upload.addEventListener('error', () => {
      this.state = 'failed';
      this.onErrorCallback();
    });

    this.xhr.upload.addEventListener('abort', () => {
      this.state = 'stopped';
      this.onAbortCallback();
    });

    this.xhr.open('POST', 'http://localhost:3001/upload/' + this.token, true);
    this.xhr.send(this.file.slice(this.bytesPreviouslySent));

    this.state = 'started';
    this.onStartCallback();
  }

  currentState() {
    return this.state;
  }

  stop() {
    if (this.state === 'started') {
      this.xhr.abort();
    }
  }

  size() {
    return this.file.size;
  }

  name() {
    return this.file.name;
  }
}

export default FileData;

