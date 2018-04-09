class Config {
  static root = () => {
    return window.apiRoot;
  }

  static get = (k) => {
    return localStorage.getItem(k);
  }

  static set = (k, v) => {
    localStorage.setItem(k, v);
  }
}

export default Config;
