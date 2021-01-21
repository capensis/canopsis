import { saveAs } from 'file-saver';

/**
 * Save file as json
 *
 * @param {Object} data
 * @param {string} name
 * @param {string} mime
 */
export const saveJsonFile = (data, name, { mime = 'application/json;charset=utf-8' } = {}) => {
  const blob = new Blob([JSON.stringify(data)], { type: mime });
  saveAs(blob, `${name}.json`);
};

/**
 * Save file as csv
 *
 * @param {string} string
 * @param {string} name
 * @param {string} mime
 */
export const saveCsvFile = (string, name, { mime = 'text/csv;charset=utf-8;' } = {}) => {
  const blob = new Blob([string], { type: mime });
  saveAs(blob, `${name}.csv`);
};
