import Rusha from 'rusha';

class HashData {
  constructor(key, data) {
    this.data = data;
    this.onHashCompleteCallback = null;
    this.key = key;
  }

  onHashComplete(callback) {
    this.onHashCompleteCallback = callback;
  }

  hash() {
    var worker = Rusha.createWorker();

    worker.onmessage = (e) => {
      var h = e.data.hash;
      this.onHashCompleteCallback(h);
    };

    worker.postMessage({id: this.key, data: this.data});
  };
}

export default HashData;
