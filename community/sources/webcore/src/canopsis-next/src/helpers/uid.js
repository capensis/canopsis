import baseUid from 'uid';

export const uid = (prefix, suffix) => {
  let result = baseUid();

  if (prefix) {
    result = `${prefix}-${result}`;
  }

  if (suffix) {
    result += `-${suffix}`;
  }

  return result;
};
