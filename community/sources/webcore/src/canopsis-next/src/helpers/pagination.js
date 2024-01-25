import { PAGINATION_LIMIT } from '@/config';

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
