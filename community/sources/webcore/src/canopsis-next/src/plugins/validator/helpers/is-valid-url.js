/**
 * Function for check url is valid
 *
 * @param {string} string
 * @return {boolean}
 */
export const isValidUrl = (string) => {
  try {
    const url = new URL(string.replace(/{{(.+)}}/g, ''));

    return /^(https?:\/\/)/.test(url.href);
  } catch (err) {
    return false;
  }
};
