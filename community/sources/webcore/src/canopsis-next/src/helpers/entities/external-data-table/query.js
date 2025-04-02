import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS } from '@/constants';

/**
 * This function converts widget with type 'AlarmsList' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export const convertExternalDataTableWidgetToQuery = (widget = {}) => {
  const { sort } = widget.parameters ?? {};

  const query = {
    page: 1,
    itemsPerPage: PAGINATION_LIMIT,
    sortBy: [],
    sortDesc: [],
  };

  if (sort?.column && sort?.order) {
    query.sortBy = [sort.column];
    query.sortDesc = [sort.order === SORT_ORDERS.desc];
  }

  return query;
};

/**
 * This function converts userPreference with widget type 'AlarmsList' to query Object
 *
 * @param {Object} [userPreference]
 * @returns {{}}
 */
export const convertExternalDataTableUserPreferenceToQuery = ({ content } = {}) => (
  content.itemsPerPage ? { itemsPerPage: content.itemsPerPage } : {}
);
