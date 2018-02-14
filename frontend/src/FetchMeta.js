class FetchMeta {
  constructor() {
    this.onFoundCallback = null;
    this.onNotFonudCallback = null;
  }

  onFound(callback) {
    this.onFoundCallback = callback;
  }

  onNotFound(callback) {
    this.onNotFoundCallback = callback;
  }

  fetch(hash) {
    var pxhr = new XMLHttpRequest();
    pxhr.addEventListener('readystatechange', (e) => {
      if (pxhr.readyState !== 4) {
        return
      }

      var text = e.target.responseText;
      if (text === undefined) {
        return;
      }

      var response = JSON.parse(text);

      if (pxhr.status === 200) {
        this.onFoundCallback(response);
      }

      if (pxhr.status === 404) {
        this.onNotFoundCallback(response);
      }
    });

    pxhr.open('GET', 'http://localhost:3001/meta/' + hash, true);
  }
}

export default FetchMeta;
