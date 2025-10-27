import { ref } from 'vue';

import { convertQueryToRequest } from '@/helpers/query';

import { usePendingHandler } from './pending';
import { useLocalQuery } from './local-query';
import { useQueryOptions } from './options';

/**
 * Custom hook that combines the functionalities of tracking pending state of an asynchronous operation and managing a
 * local query state.
 *
 * This hook integrates `usePendingHandler` to manage the pending state of an asynchronous operation and `useLocalQuery`
 * to handle local query state updates.
 * It triggers the asynchronous operation whenever the query state updates.
 *
 * @param {Object} [options] - Configuration options for the hook.
 * @param {Object} [options.initialQuery] - The initial state of the query object.
 * @param {boolean} [options.initialPending] - The initial pending state for the asynchronous operation.
 * @param {Function} [options.comparator] - Function used to compare the old and new query values.
 * @param {Function} [options.fetchHandler] - The asynchronous function that will be executed when the query updates.
 * @returns {Object} An object containing the combined functionalities of pending state and query management.
 */
export const usePendingWithLocalQuery = ({
  initialQuery,
  initialPending,
  comparator,
  fetchHandler,
} = {}) => {
  const {
    pending,
    handler: fetchHandlerWithPending,
  } = usePendingHandler(fetchHandler, initialPending);

  const queryData = useLocalQuery({
    initialQuery,
    comparator,
    onUpdate: fetchHandlerWithPending,
  });

  const fetchHandlerWithQuery = (requestQuery = queryData.query.value) => fetchHandlerWithPending(requestQuery);

  return {
    ...queryData,

    pending,

    fetchHandlerWithQuery,
  };
};

/**
 * Custom hook to manage local query state with options handling.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Object} [options.initialQuery] - The initial state of the query object.
 * @param {Function} [options.comparator] - Function used to compare the old and new query values.
 * @param {Function} [options.onUpdate] - Callback function triggered when the query updates.
 * @returns {Object} An object containing query state, options, and update functions.
 */
export const useLocalQueryWithOptions = ({ initialQuery, comparator, onUpdate }) => {
  const queryData = useLocalQuery({
    initialQuery,
    comparator,
    onUpdate,
  });

  const { options, updateOptions } = useQueryOptions(queryData.query, queryData.updateQuery);

  return {
    ...queryData,
    options,
    updateOptions,
  };
};

/**
 * Custom hook to fetch a list of data with local query and options management.
 *
 * This hook combines local query state management with options (such as pagination and sorting),
 * and triggers the provided fetchListHandler callback with the query converted to request params.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Function} options.fetchListHandler - Async function that fetches data based on provided parameters.
 * @param {Object} [options.initialQuery] - The initial state of the query object.
 * @returns {Object} An object containing query state, options, and update functions.
 */
export const useFetchListWithOptions = ({ fetchListHandler, initialQuery }) => useLocalQueryWithOptions({
  initialQuery,
  onUpdate: fetchQuery => fetchListHandler({ params: convertQueryToRequest(fetchQuery) }),
});

/**
 * Custom hook to fetch a list of data with query options management
 *
 * This hook combines data fetching with query state management and options handling.
 * It provides reactive references for the fetched data, metadata, and loading state,
 * along with functions to manage the query and pagination/sorting options.
 *
 * @param {Object} options - Configuration options for the hook
 * @param {Function} options.fetchListHandler - Async function that fetches data based on provided parameters
 * @param {Object} [options.initialQuery] - Initial query state
 * @returns {Object}
 */
export const useFetchListWithoutStoreWithOptions = ({ fetchListHandler, initialQuery }) => {
  const data = ref([]);
  const meta = ref({});

  const {
    pending,
    query,
    updateQuery,
    resetQuery,
    fetchHandlerWithQuery: fetchList,
  } = usePendingWithLocalQuery({
    initialQuery,
    fetchHandler: async (fetchQuery) => {
      const response = await fetchListHandler({ params: convertQueryToRequest(fetchQuery) });

      data.value = response.data;
      meta.value = response.meta;
    },
  });

  const { options, updateOptions } = useQueryOptions(query, updateQuery);

  return {
    data,
    meta,
    pending,
    fetchList,

    query,
    updateQuery,
    resetQuery,

    options,
    updateOptions,
  };
};
