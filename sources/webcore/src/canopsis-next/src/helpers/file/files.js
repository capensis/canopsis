import { saveAs } from 'file-saver';

/**
 * Save data in file
 *
 * @param {Blob} blob
 * @param {string} name
 * @param {string} extension
 */
export const saveFile = (blob, name, extension) => saveAs(blob, `${name}.${extension}`);

/**
 * Save file as json
 *
 * @param {Object} data
 * @param {string} name
 */
export const saveJsonFile = (data, name) => {
  const blob = new Blob([JSON.stringify(data)], { type: 'application/json;charset=utf-8' });

  return saveFile(blob, name, 'json');
};

/**
 * Save file as csv
 *
 * @param {string} string
 * @param {string} name
 */
export const saveCsvFile = (string, name) => {
  const blob = new Blob([string], { type: 'text/csv;charset=utf-8;' });

  return saveFile(blob, name, 'csv');
};
