import uid from 'uid';

export default (prefix, suffix) => {
  let result = uid();

  if (prefix) {
    result = `${prefix}-${result}`;
  }

  if (suffix) {
    result += `-${suffix}`;
  }

  return result;
};
