const Config = {
  url: 'http://localhost:3001',
  postFileUrl: (hash) => (`${Config.url}/api/file/${hash}`),
  getMetaUrl: (hash) => (`${Config.url}/api/meta/by-hash/${hash}`),
  postMetaUrl: (h) => (`${Config.url}/api/meta`),
  getFileBySlugUrl: (slug) => (`${Config.url}/api/file/by-slug/${slug}`),
  getFilesListUrl: (page) => (`${Config.url}/api/files?page=${page}`),
};

export default Config;
