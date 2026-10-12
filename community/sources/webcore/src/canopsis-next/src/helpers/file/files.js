import { saveAs } from 'file-saver';

/**
 * Save data in file
 *
 * @param {Blob} blob
 * @param {string} name
 */
export const saveFile = (blob, name) => saveAs(blob, name);

/**
 * Save file as json
 *
 * @param {Object|string|number} data
 * @param {string} name
 */
export const saveJsonFile = (data, name) => {
  const blob = new Blob([JSON.stringify(data)], { type: 'application/json;charset=utf-8' });

  return saveFile(blob, `${name}.json`);
};

/**
 * Save file as csv
 *
 * @param {string} string
 * @param {string} name
 */
export const saveCsvFile = (string, name) => {
  const blob = new Blob([string], { type: 'text/csv;charset=utf-8;' });

  return saveFile(blob, `${name}.csv`);
};

/**
 * Save file as txt
 *
 * @param {string} text
 * @param {string} name
 */
export const saveTextFile = (text, name) => {
  const blob = new Blob([text], { type: 'text/plain;charset=utf-8;' });

  return saveFile(blob, `${name}.txt`);
};
