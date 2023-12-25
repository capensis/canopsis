import { PAGINATION_LIMIT } from '@/config';

/**
 * Get page for updated recordsPerPage
 *
 * @param {number} newRecordsPerPage
 * @param {number} [oldRecordsPerPage = PAGINATION_LIMIT]
 * @param {number} [oldPage = 1]
 * @return {number}
 */
export const getPageForNewRecordsPerPage = (
  newRecordsPerPage,
  oldRecordsPerPage = PAGINATION_LIMIT,
  oldPage = 1,
) => (
  Math.ceil(Math.min(oldRecordsPerPage * (oldPage - 1) + 1) / newRecordsPerPage)
);
