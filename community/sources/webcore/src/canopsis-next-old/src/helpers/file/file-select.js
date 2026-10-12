/**
 * Returns a promise that resolves with the result of the FileReader once the file is loaded.
 *
 * @param {FileReader} reader - The FileReader instance used to read the file.
 * @returns {Promise<string|ArrayBuffer>} A promise that resolves with the result of the FileReader.
 */
export const getFileReaderResult = reader => new Promise((resolve, reject) => {
  reader.addEventListener('load', e => resolve(e.target.result));
  reader.addEventListener('error', reject);
});

/**
 * Reads the content of a file as text.
 *
 * @param {File} file - The file to be read.
 * @returns {Promise<string>} A promise that resolves with the text content of the file.
 */
export const getFileTextContent = (file) => {
  const reader = new FileReader();

  reader.readAsText(file, 'UTF-8');

  return getFileReaderResult(reader);
};

/**
 * Reads the content of a file as a data URL.
 *
 * @param {File} file - The file to be read.
 * @returns {Promise<string>} A promise that resolves with the data URL of the file.
 */
export const getFileDataUrlContent = (file) => {
  const reader = new FileReader();

  reader.readAsDataURL(file);

  return getFileReaderResult(reader);
};
