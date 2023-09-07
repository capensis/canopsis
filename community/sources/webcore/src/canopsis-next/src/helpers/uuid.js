import { v4 } from 'uuid';

export const uuid = (prefix, suffix) => {
  let result = v4();

  if (prefix) {
    result = `${prefix}_${result}`;
  }

  if (suffix) {
    result += `_${suffix}`;
  }

  return result;
};
