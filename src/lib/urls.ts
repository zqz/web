const isProduction = import.meta.env.PROD;
const productionUrl = 'https://api.zqz.ca';
const devUrl = 'http://localhost:3001/api';
const url = isProduction ? productionUrl : devUrl;

const URLs = {
  url: url,
  postFileUrl: (hash: string) => (`${url}/file/${hash}`),
  getMetaUrl: (hash: string) => (`${url}/meta/by-hash/${hash}`),
  postMetaUrl: () => (`${url}/meta`),
  getFileBySlugUrl: (slug: string) => (
    isProduction ?
    `https://x.zqz.ca/${slug}` :
    `${url}/file/by-slug/${slug}`
  ),
  getFilesListUrl: (page: string) => (`${url}/files?page=${page}`),
  thumbnailUrl: (slug: string) => (`https://thumbnails.zqz.ca/${slug}`)
};

export default URLs;
