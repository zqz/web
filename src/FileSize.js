const bytes = (value) => {
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    if (value === 0) return '0';
    var i = parseInt(Math.floor(Math.log(value) / Math.log(1024)), 10);
    var v = value / Math.pow(1024, i);
    var n = '';

    if (i === 0) {
      n = Math.round(v, 2);
    } else {
      n = v.toFixed(2);
    }

    return n + ' ' + sizes[i];
  }

export default bytes;

