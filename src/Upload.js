const Upload = (file) => {
  function upload() {
    console.log('upload');
    if (this.onStart !== null) {
      this.onStart('starting');
    }
  }

  return {
    upload
  }
}

export default Upload;
