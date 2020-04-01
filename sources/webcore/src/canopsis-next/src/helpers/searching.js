import escapeStringRegexp from 'escape-string-regexp';

export function replaceTextNotInQuotes(str, text, value) {
  const escapedText = escapeStringRegexp(text);
  const regexp = new RegExp(`(?!\\B"[^"]*)${escapedText}(?![^"]*"\\B)|(?!\\\\B'[^']*)${escapedText}(?![^']*'\\\\B)`, 'gu');

  return str.replace(regexp, value);
}

/**
 * Calculate duplicate names count
 * @param {Array} list
 * @param {Object} item
 * @return {Number}
 * Example: getDuplicateCountItems([{ name: 'test' }, { name: 'test' }, { name: 'other' }], { name: 'test' }) -> 2
 */
export const getDuplicateCountItems = (list, item) => {
  const nameRegexp = new RegExp(`^${item.name}(\\s\\(\\d+\\))?$`);

  return list.reduce((count, { name }) => {
    const isDuplicate = nameRegexp.test(name);

    return isDuplicate ? count + 1 : count;
  }, 0);
};
