class Config {
  static root = () => {
    return window.location.protocol + '//' + window.location.hostname + ':3001';
  }

  static get = (k) => {
    return localStorage.getItem(k);
  }

  static set = (k, v) => {
    localStorage.setItem(k, v);
  }
}

export default Config;
