import { orderBy } from 'lodash';

import { replaceTextNotInQuotes } from './quotes';

/**
 * Flattens a nested array of column objects into a single-level array.
 *
 * @param {WidgetColumn[]} [columns = []] - An array of column objects to be flattened. Each column
 *                                  object may contain nested child columns.
 * @param {string} [childrenKey = 'items'] - The key under which child columns are stored in
 *                                         each column object. Defaults to 'items'.
 * @returns {WidgetColumn[]} A flattened array of column objects.
 */
export const flattenSearchColumns = (columns = [], childrenKey = 'items') => columns.reduce((acc, column) => {
  acc.push(column);

  if (column[childrenKey]?.length) {
    acc.push(...flattenSearchColumns(column[childrenKey]));
  }

  return acc;
}, []);

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

  const orderedColumns = orderBy(
    flattenSearchColumns(columns),
    ({ selectorText = '', text = '' }) => selectorText.length || text.length,
    ['desc'],
  );

  return orderedColumns.reduce(
    (acc, { selectorText, text, value }) => replaceTextNotInQuotes(acc, selectorText || text, value),
    search.replace(/^-(\s*)/, ''),
  );
};
