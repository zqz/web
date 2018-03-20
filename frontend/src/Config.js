class Config {
  static root = () => {
    return window.location.protocol + '//' + window.location.hostname + ':3001';
  }
}

export default Config;
