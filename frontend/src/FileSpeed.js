import FileSize from './FileSize';

class FileSpeed {
  constructor(size) {
    this.startTime = null;
    this.size = 0;
    this.duration = 0;
    this.entries = [];
  }

  add(bytes, time) {
    if (time === undefined) {
      var d = new Date();
      time = d.getTime();
    }

    if (this.startTime === null) {
      this.startTime = time;
      console.log("start time:" + this.startTime);
    }

    this.entries.push([time, bytes]);

    console.log(this.entries);
  }

  speed() {
    var s = null;
    var total = 0;
    var duration = null;

    if (this.entries.length === 0) {
      return '';
    }

    for(var i = 0; i < this.entries.length; i++) {
      var e = this.entries[i];
      var t = e[0];
      var b = e[1];

      if (s === null) {
        s = t;
      }

      var st = t - s;
      total = b;
      duration = st;
    }

    console.log(total, duration);

    return FileSize.ToSize((total * 1024) / duration) + '/s';
  }
}

export default FileSpeed;
