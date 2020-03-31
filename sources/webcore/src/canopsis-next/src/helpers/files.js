import { saveAs } from 'file-saver';

export const saveJsonFile = (data, name, { mime = 'application/json;charset=utf-8' } = {}) => {
  const blob = new Blob([JSON.stringify(data)], { type: mime });
  saveAs(blob, `${name}.json`);
};
