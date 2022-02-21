import { isEmpty } from 'lodash';

/**
 * Function for check picker value is valid
 *
 * @param {[] | string | number} value
 * @param {Function} [preparer]
 * @return {boolean}
 */
export const isValidPicker = (value, { preparer }) => {
  try {
    if (!isEmpty(value)) {
      return Boolean(preparer ? preparer(value) : value);
    }

    return true;
  } catch (err) {
    return false;
  }
};
