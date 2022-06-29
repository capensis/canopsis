/**
 * Get collection comparator for the sort by special field
 *
 * @param {string} field
 * @returns {function(...[*]=)}
 */
export function getCollectionComparator(field) {
  return (a, b) => {
    if (a[field] > b[field]) {
      return 1;
    }

    if (a[field] < b[field]) {
      return -1;
    }

    return 0;
  };
}
