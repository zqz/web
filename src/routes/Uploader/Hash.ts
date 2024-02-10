const Hash = (file: File, callback: any) => {
  function hex(buf) {
    const hashArray = Array.from(new Uint8Array(buf));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  }

  async function onBuffer(b) {
    const buf = await crypto.subtle.digest('SHA-1', b);
    const hash = hex(buf);
    callback(hash);
  }

  if (file.data) {
    file.data.arrayBuffer().then(onBuffer);
  }
}

export default Hash;
