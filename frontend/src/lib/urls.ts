const isProduction = import.meta.env.PROD;
const productionUrl = 'https://api.zqz.ca';
const devUrl = 'http://localhost:3001/api';
const url = isProduction ? productionUrl : devUrl;

function getFileBySlugUrl(slug: string) {
  if (isProduction) {
    return `https://x.zqz.ca/${slug}`;
  }

  return `${url}/file/by-slug/${slug}`;
} 

export const URLs = {
  url: url,
  postFileUrl: (hash: string) => `${url}/file/${hash}`,
  getMetaUrl: (hash: string) => `${url}/meta/by-hash/${hash}`,
  postMetaUrl: () => `${url}/meta`,
  getFileBySlugUrl: getFileBySlugUrl,
  getFilesListUrl: (page: number) => (`${url}/files?page=${page}`),
  thumbnailUrl: (slug: string) => isProduction ? `https://thumbnails.zqz.ca/${slug}` : getFileBySlugUrl(slug)
};
