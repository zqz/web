import type { Meta } from "$lib/types";

export function fileIsImage(file: Meta) : boolean {
  const t = extractMediaType(file);
  return isImage(t);
}

function isImage(mediaType: string) : boolean {
  return mediaType === 'image';
}

function extractMediaType(file: Meta) {
  const parts = file.type.split('/');
  return parts[0];
}
