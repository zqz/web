class Meta {
  constructor(meta) {
    this.alias = meta.alias;
    this.name = meta.name;
    this.hash = meta.hash;
    this.slug = meta.slug;
    this.type = meta.type;
    this.path = meta.path;
    this.size = meta.size;
    this.date = meta.date;
    this.bytesReceived = meta.bytesReceived;
  }

  obj = () => {
    return {
      alias: this.alias,
      name: this.name,
      hash: this.hash,
      slug: this.slug,
      type: this.type,
      path: this.path,
      size: this.size,
      date: this.date,
      bytesReceived: this.bytesReceived
    }
  }

  data() {
    return JSON.stringify(this.meta.obj());
  }

  post() {
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

      this._callback(response);
    });

    pxhr.open('POST', 'http://localhost:3001/data/meta', true);
    pxhr.send(this.data());
  }
}

export default Meta;
