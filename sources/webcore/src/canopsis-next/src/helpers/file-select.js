export function getFileReaderResult(reader) {
  return new Promise((resolve, reject) => {
    reader.addEventListener('load', e => resolve(e.target.result));
    reader.addEventListener('error', reject);
  });
}

export function getFileTextContent(file) {
  const reader = new FileReader();

  reader.readAsText(file, 'UTF-8');

  return getFileReaderResult(reader);
}

export function getFileDataUrlContent(file) {
  const reader = new FileReader();

  reader.readAsDataURL(file);

  return getFileReaderResult(reader);
}
