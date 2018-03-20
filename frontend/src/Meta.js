import Config from './Config';

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
}

export default Meta;
