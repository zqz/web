import Config from "./Config";

class FetchMeta {
  constructor() {
    this.onFoundCallback = null;
    this.onNotFoundCallback = null;
  }

  onFound(callback) {
    this.onFoundCallback = callback;
  }

  onNotFound(callback) {
    this.onNotFoundCallback = callback;
  }

  fetch(hash) {
    if (hash === undefined) {
      return;
    }

    var pxhr = new XMLHttpRequest();
    pxhr.addEventListener("readystatechange", (e) => {
      if (pxhr.readyState !== 4) {
        return
      }

      var text = e.target.responseText;
      if (text === undefined || text === "") {
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

    pxhr.open("GET", Config.root() + "/meta/" + hash, true);
    pxhr.send();
  }
}

export default FetchMeta;
