const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];

const bytes = (value: number) : string => {
  if (value === 0 || isNaN(value)) return '0';
  const i = Math.floor(Math.log(value) / Math.log(1024));
  const v = value / Math.pow(1024, i);
  let n = v;

  if (i === 0) {
    n = Math.round(v);
  }

  return `${n.toFixed(2)} ${sizes[i]}`;
}

export default bytes;
