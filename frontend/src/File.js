// var shajs = require('sha.js')
import Meta from './Meta.js';
import FetchMeta from './FetchMeta.js';
import PostMeta from './PostMeta.js';
import HashData from './HashData.js';
import PostData from './PostData.js';

class FileData {
  constructor(data) {
    // this.state = null;

    this._meta = new Meta(data);
    this._data = data;
    this._key = Math.random().toString(36).substring(7);

    this._fetchMeta = new FetchMeta();
    this._postMeta = new PostMeta();
    this._hashData = new HashData(this._key, data);
    this._postData = new PostData(data);

    this._postMeta.onResponse((r) => {
      console.log('on post meta: ', r);
      this._meta.bytesReceived = r.bytes_received;
      this._postData.post(this._meta.hash, r.bytes_received);
    });

    this.onExistsCallback = null;
  }

  meta() {
    return this._meta;
  }

  key() {
    return this._key;
  }

  onProgress(callback) {
    this._postData.onProgress(callback);
  }

  onError(callback) {
    this._postData.onError(callback);
  }

  onStart(callback) {
    this._postData.onStart(callback);
  }

  onAbort(callback) {
    this._postData.onAbort(callback);
  }

  onLoad(callback) {
    this._postData.onLoad(callback);
  }

  onResponse(callback) {
    this._postData.onResponse(callback);
  }

  onHash(callback) {
    this._hashData.onHashComplete((hash) => {
      this._meta.hash = hash;
      // this.check();
      callback(hash)
    });
  }

  // not sure
  onExists(callback) {
    this.onExistsCallback = callback;
  }

  check() {
    var fm = new FetchMeta(this.meta.hash);
    fm.onFound((json) => {

    });

    fm.onNotFound((json) => {
    });

    fm.fetch();
  }

  start() {
    this._postMeta.post(this._meta);
  }

  stop() {
    this._postData.abort();
  }

  hash() {
    this._hashData.hash()
  }
}

export default FileData;

