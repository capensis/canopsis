/**
 * Function for check value is unique
 *
 * @param {[] | string | number} inputValue
 * @param {[]} values
 * @param {string | number} initialValue
 * @return {boolean}
 */
export const isUniqueValue = (inputValue, { values, initialValue }) => {
  if (!inputValue) {
    return false;
  }

  if (Array.isArray(inputValue)) {
    if (!inputValue.length) {
      return true;
    }

    return !inputValue.some(value => values.includes(value.toLowerCase()));
  }

  if (initialValue === inputValue) {
    return true;
  }

  return !values.includes(inputValue.toLowerCase());
};
