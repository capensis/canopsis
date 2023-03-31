import { isObject, isUndefined, isEmpty } from 'lodash';

export const differenceKeys = (object, base) => Object.keys(object).reduce((acc, key) => {
  if (isUndefined(base[key])) {
    acc[key] = true;
  } else if (isObject(object[key]) && isObject(base[key])) {
    const localDiff = differenceKeys(object[key], base[key]);

    if (!isEmpty(localDiff)) {
      acc[key] = localDiff;
    }
  }

  return acc;
}, {});
