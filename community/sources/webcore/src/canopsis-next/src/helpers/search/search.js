import { replaceTextNotInQuotes } from './quotes';

/**
 * Prepares a search string for submission by replacing column text with corresponding values if the search string
 * starts with a hyphen. If the search string does not start with a hyphen, it is returned as is.
 *
 * @param {string} search - The search string to be processed.
 * @param {WidgetColumn[]} [columns = []] - Optional array of objects where each object contains `text` to be replaced
 * and its `value`.
 * @returns {string} - The processed search string with replacements made, or the original search string if no hyphen
 * prefix is present.
 */
export const prepareSearchForSubmit = (search = '', columns = []) => {
  if (!search.startsWith('-')) {
    return search;
  }

  return columns.reduce(
    (acc, { text, value }) => replaceTextNotInQuotes(acc, text, value),
    search.replace(/^-(\s*)/, ''),
  );
};
