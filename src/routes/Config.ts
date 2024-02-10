const isProduction = import.meta.env.PROD;

const Config = {
  url: isProduction ? 'https://api.zqz.ca' : 'http://localhost:3001/api',
  postFileUrl: (hash: string) => (`${Config.url}/file/${hash}`),
  getMetaUrl: (hash: string) => (`${Config.url}/meta/by-hash/${hash}`),
  postMetaUrl: (_: string) => (`${Config.url}/meta`),
  getFileBySlugUrl: (slug: string) => (
    isProduction ?
    `https://x.zqz.ca/${slug}` :
    `${Config.url}/file/by-slug/${slug}`
  ),
  getFilesListUrl: (page: string) => (`${Config.url}/files?page=${page}`),
  thumbnailUrl: (slug: string) => (`https://thumbnails.zqz.ca/${slug}`)
};

export default Config;
