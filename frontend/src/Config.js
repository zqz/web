class Config {
  static root = () => {
    return window.apiRoot;
  }
  static cdnroot = () => {
    return window.cdnRoot;
  }

  static get = (k) => {
    return localStorage.getItem(k);
  }

  static set = (k, v) => {
    localStorage.setItem(k, v);
  }

  static toggleDark = () => {
    if (Config.get('dark') === 'true') {
      document.body.classList.add('Dark');
    } else {
      document.body.classList.remove('Dark');
    }
  }
}

export default Config;
