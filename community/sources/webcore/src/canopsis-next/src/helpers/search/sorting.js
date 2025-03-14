import { uniq } from 'lodash';

/**
 * @typedef {Object} Search
 * @property {string} search
 * @property {boolean} pinned
 */

/**
 * Immutable sorting for pinned searches
 *
 * @param {Search[]} searches
 * @param {string} activeSearch
 * @param {string} [key = 'search']
 * @returns {Search[]}
 */
export const sortPinnedSearches = (searches, activeSearch, key = 'search') => uniq(searches).sort((a, b) => {
  const aFactor = Number(a[key] === activeSearch) + (a.pinned ? 2 : 0);
  const bFactor = Number(b[key] === activeSearch) + (b.pinned ? 2 : 0);

  return bFactor - aFactor;
});
