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
      console.log('hash undefined in fetchmeta');
      return;
    }

    console.log('hii1');

    var pxhr = new XMLHttpRequest();
    pxhr.addEventListener('readystatechange', (e) => {
      if (pxhr.readyState !== 4) {
        console.log(pxhr.readyState);
        return
      }

      console.log('hi');
      var text = e.target.responseText;
      if (text === undefined) {
        return;
      }

      var response = JSON.parse(text);

      if (pxhr.status === 200) {

        console.log('found');
        this.onFoundCallback(response);
      }

      if (pxhr.status === 404) {
        console.log('notfound');
        this.onNotFoundCallback(response);
      }
    });

    pxhr.open('GET', 'http://localhost:3001/meta/' + hash, true);
    pxhr.send();
  }
}

export default FetchMeta;
