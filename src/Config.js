const isProduction = document.location.hostname !== 'localhost';

const Config = {
  url: isProduction ? 'https://api.zqz.ca' : 'http://localhost:3001/api',
  postFileUrl: (hash) => (`${Config.url}/file/${hash}`),
  getMetaUrl: (hash) => (`${Config.url}/meta/by-hash/${hash}`),
  postMetaUrl: (h) => (`${Config.url}/meta`),
  getFileBySlugUrl: (slug) => (
    isProduction ?
    `https://x.zqz.ca/${slug}` :
    `${Config.url}/file/by-slug/${slug}`
  ),
  getFilesListUrl: (page) => (`${Config.url}/files?page=${page}`),
  thumbnailUrl: (slug) => (`https://thumbnails.zqz.ca/${slug}`)
};

export default Config;
