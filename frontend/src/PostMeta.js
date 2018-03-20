import Config from './Config';

class PostMeta {
  constructor() {
    this.onResponseCallback = null;
    this.onFinishedCallback = null;

    this.onResponse = this.onResponse.bind(this);
    this.onFinished = this.onFinished.bind(this);
  }

  onResponse(callback) {
    this.onResponseCallback = callback;
  }

  onFinished (callback) {
    this.onFinishedCallback = callback;
  }

  data() {
  }

  post(meta) {
    if (meta === undefined) {
      console.log("meta is a undefined");
      return;
    }

    var pxhr = new XMLHttpRequest();
    pxhr.addEventListener('readystatechange', (e) => {
      if (pxhr.readyState !== 4) {
        return;
      }

      var text = e.target.responseText;

      if (text === undefined || text === '') {
        return;
      }

      var response = JSON.parse(text);

      this.onResponseCallback(response);
    });

    var data = JSON.stringify(meta.obj());
    pxhr.open('POST', Config.root() + '/meta', true);
    pxhr.send(data);
  }
}

export default PostMeta;
