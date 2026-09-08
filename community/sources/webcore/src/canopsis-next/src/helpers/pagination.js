import { PAGINATION_LIMIT } from '@/config';

/**
 * Returns the greatest valid 1-based page number when paginating `itemsCount` items with fixed page size.
 *
 * @param {number} itemsCount - Total number of rows in the dataset.
 * @param {number} [itemsPerPage] - Rows per page; coerced via app default when missing or zero.
 *
 * @returns {number}
 */
export const getMaxPageForItemsCount = (itemsCount, itemsPerPage) => {
  const perPage = itemsPerPage || PAGINATION_LIMIT || 10;

  return Math.max(1, Math.ceil(itemsCount / perPage) || 1);
};

/**
 * Get page for updated itemsPerPage
 *
 * @param {number} newItemsPerPage
 * @param {number} [oldItemsPerPage = PAGINATION_LIMIT]
 * @param {number} [oldPage = 1]
 * @return {number}
 */
export const getPageForNewItemsPerPage = (
  newItemsPerPage,
  oldItemsPerPage = PAGINATION_LIMIT,
  oldPage = 1,
) => (
  Math.ceil(Math.min(oldItemsPerPage * (oldPage - 1) + 1) / newItemsPerPage)
);
