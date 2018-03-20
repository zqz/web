import Config from './Config';

class PostData {
  constructor(data) {
    this.onProgressCallback = null;
    this.onStartCallback = null;
    this.onAbortCallback = null;
    this.onErrorCallback = null;
    this.onFinishCallback = null;
    this.onResponseCallback = null;

    this.data = data;
    this.xhr = null;
  }

  onProgress(callback) {
    this.onProgressCallback = callback;
  };

  onError(callback) {
    this.onErrorCallback = callback;
  };

  onStart(callback) {
    this.onStartCallback = callback;
  };

  onAbort(callback) {
    this.onAbortCallback = callback;
  };

  onLoad(callback) {
    this.onLoadCallback = callback;
  };

  onResponse(callback) {
    this.onResponseCallback = callback;
  };

  post(hash, offset) {
    if (hash === undefined) {
      console.log('hash undefined');
      return;
    }

    if (offset === undefined) {
      console.log('offset undefined');
      return;
    }

    this.xhr = new XMLHttpRequest();
    this.xhr.upload.addEventListener('progress', (e) => {
      this.onProgressCallback({
        loaded: offset + e.loaded,
        total: this.data.size
      });
    });

    this.xhr.upload.addEventListener('load', () => {
      this.onLoadCallback();
    });

    this.xhr.addEventListener('readystatechange', (e) => {
      if (this.xhr.readyState !== 4) {
        return;
      }
      var text = e.target.responseText;

      if (text === undefined || text === null || text === '') {
        return;
      }

      var data = JSON.parse(text);
      console.log('postdata', data);
      this.onResponseCallback(data);
    });

    this.xhr.upload.addEventListener('error', () => {
      this.onErrorCallback();
    });

    this.xhr.upload.addEventListener('abort', () => {
      this.onAbortCallback();
    });

    this.xhr.open('POST', Config.root() + '/data/' + hash, true);
    this.xhr.send(this.data.slice(offset));
    this.onStartCallback();
  }

  abort() {
    this.xhr.abort();
  };
}

export default PostData;
