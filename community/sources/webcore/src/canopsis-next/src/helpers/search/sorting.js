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
 * @returns {Search[]}
 */
export const immutableSortPinnedSearches = (searches, activeSearch) => [...searches].sort((a, b) => {
  const aFactor = Number(a.search === activeSearch) + Number(a.pinned);
  const bFactor = Number(b.search === activeSearch) + Number(b.pinned);

  return bFactor - aFactor;
});
