function hex(buf: ArrayBuffer) : string {
  const arr = Array.from(new Uint8Array(buf));
  return arr.map(b => b.toString(16).padStart(2, '0')).join('');
}

export const hashFile = (file: globalThis.File, callback: (h: string) => void) => {
  async function onBuffer(b: BufferSource) {
    const buf = await crypto.subtle.digest('SHA-1', b);
    const hash = hex(buf);
    callback(hash);
  }

  if (file) {
    file.arrayBuffer().then(onBuffer);
  }
}
