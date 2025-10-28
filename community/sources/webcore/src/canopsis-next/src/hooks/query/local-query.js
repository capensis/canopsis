import { isEqual, upperFirst, camelCase } from 'lodash';
import { ref, unref, set } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

/**
 * Custom hook for managing local query state with Vue.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Object} options.initialQuery - The initial state of the query object.
 * @param {Function} options.onUpdate - Callback function that is called when the query is updated.
 * @param {Function} options.comparator - Function used to compare the old and new query values.
 * @returns {Object} An object containing methods to manipulate the query state and the current query state itself.
 *
 * @example
 * // Usage of useLocalQuery within a Vue component
 * import { useLocalQuery } from `./path/to/useLocalQuery`;
 *
 * export default {
 *   setup() {
 *     const { query, updateQueryPage, updateQueryLimit } = useLocalQuery({
 *       initialQuery: { page: 1, limit: 10 },
 *       onUpdate: () => console.log(`Query updated`),
 *       comparator: (a, b) => JSON.stringify(a) === JSON.stringify(b)
 *     });
 *
 *     // Update the page number
 *     updateQueryPage(2);
 *
 *     // Update the limit
 *     updateQueryLimit(20);
 *
 *     return {
 *       query
 *     };
 *   }
 * };
 */
export const useLocalQuery = ({
  initialQuery = { page: 1, itemsPerPage: PAGINATION_LIMIT },
  onUpdate,
  comparator = isEqual,
} = {}) => {
  const query = ref({ ...unref(initialQuery) });

  /**
   * Updates the entire query object with a new query.
   * Triggers the onUpdate callback if the new query differs from the old one according to the comparator.
   *
   * @param {Object} newQuery - The new query object to replace the current query.
   */
  const updateQuery = (newQuery) => {
    const oldQuery = query.value;

    query.value = newQuery;

    if (onUpdate && !comparator(oldQuery, newQuery)) {
      onUpdate(newQuery);
    }
  };

  /**
   * Updates a specific field in the query object.
   * Triggers the onUpdate callback if the new value differs from the old value.
   *
   * @param {string} field - The field name to update in the query object.
   * @param {*} value - The new value for the specified field.
   */
  const updateQueryField = (field, value) => {
    const oldValue = query.value?.[field];

    set(query.value, field, value);

    if (onUpdate && !isEqual(oldValue, value)) {
      onUpdate(query.value);
    }
  };

  /**
   * Resets the query to its initial state.
   */
  const resetQuery = () => query.value = { ...unref(initialQuery) };

  /**
   * Calls the onUpdate callback with the provided query or the current query value.
   *
   * @param {Object} [handlerQuery=query.value] - The query object to pass to the onUpdate callback.
   */
  const handler = (handlerQuery = query.value) => onUpdate(handlerQuery);

  const updateQueryFieldsMethods = Object.keys(unref(initialQuery)).reduce((acc, field) => {
    acc[`updateQuery${upperFirst(camelCase(field))}`] = value => updateQueryField(field, value);

    return acc;
  }, {});

  return {
    ...updateQueryFieldsMethods,

    query,
    updateQuery,
    updateQueryField,
    resetQuery,
    handler,
  };
};
