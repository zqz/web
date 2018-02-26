// var shajs = require('sha.js')
import Meta from './Meta.js';
import FetchMeta from './FetchMeta.js';
import PostMeta from './PostMeta.js';
import HashData from './HashData.js';
import PostData from './PostData.js';
import FileSpeed from './FileSpeed.js';

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
    this._fileSpeed = new FileSpeed();

    this._postMeta.onResponse((r) => {
      console.log('on post meta: ', r);
      this._meta.bytesReceived = r.bytes_received;
      this._postData.post(this._meta.hash, r.bytes_received);
    });

    this.onExistsCallback = null;
  }

  speed() {
    return this._fileSpeed.speed();
  }

  meta() {
    return this._meta;
  }

  key() {
    return this._key;
  }

  onFound(callback) {
    this._fetchMeta.onFound((data) => {
      this._meta.bytesReceived = data.bytes_received;
      callback(data);
    });
  }

  onNotFound(callback) {
    this._fetchMeta.onNotFound(callback);
  }

  onProgress(callback) {
    this._postData.onProgress((data) => {
      this._fileSpeed.add(data.loaded - this._meta.bytesReceived);
      callback(data);
    });
  }

  onError(callback) {
    this._postData.onError(callback);
  }

  onStart(callback) {
    this._postData.onStart((d) => {
      callback(d);
    });
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

  check(hash) {
    this._fetchMeta.fetch(hash);
    // var fm = new FetchMeta(this.meta.hash);
    // fm.onFound((json) => {

    // });

    // fm.onNotFound((json) => {
    // });

    // fm.fetch();
  }

  start() {
    this._filespeed = new FileSpeed();
    if (this._meta.bytesReceived > 0) {
      this._postData.post(this._meta.hash, this._meta.bytesReceived);
    } else {
      this._postMeta.post(this._meta);
    }
  }

  stop() {
    this._postData.abort();
  }

  hash() {
    this._hashData.hash()
  }
}

export default FileData;

