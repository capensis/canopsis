/**
 * Save data in file
 *
 * @param {Blob} blob
 * @param {string} name
 */
export const saveFile = (data, filename) => {
  const blob = data instanceof Blob ? data : new Blob([data]);
  const link = document.createElement('a');

  link.href = URL.createObjectURL(blob);
  link.download = filename;
  link.click();

  URL.revokeObjectURL(link.href);
};

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
 * @param {string|Blob} data
 * @param {string} name
 */
export const saveCsvFile = (data, name) => {
  const blob = data instanceof Blob ? data : new Blob([data], { type: 'text/csv;charset=utf-8;' });

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
