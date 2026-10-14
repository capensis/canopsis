import { uid as baseUid } from 'uid';

export const uid = (prefix: string = '', suffix: string = ''): string => {
  let result = baseUid();

  if (prefix) {
    result = `${prefix}-${result}`;
  }

  if (suffix) {
    result += `-${suffix}`;
  }

  return result;
};
