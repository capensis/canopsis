/**
 * Function for check json is valid
 *
 * @param {string} json
 * @return {boolean}
 */
export const isValidJson = (json) => {
  try {
    return !!JSON.parse(json);
  } catch (err) {
    return false;
  }
};
