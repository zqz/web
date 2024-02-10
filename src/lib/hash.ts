// todo: add proper types here
//
const hashFile = (file: globalThis.File, callback: any) => {
  function hex(buf: ArrayBuffer) {
    const hashArray = Array.from(new Uint8Array(buf));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  }

  async function onBuffer(b: BufferSource) {
    const buf = await crypto.subtle.digest('SHA-1', b);
    const hash = hex(buf);
    callback(hash);
  }

  if (file) {
    file.arrayBuffer().then(onBuffer);
  }
}

export default hashFile;
