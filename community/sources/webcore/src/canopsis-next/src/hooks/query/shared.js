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
 * Custom hook to fetch a list of data without using a store, with options for query management.
 *
 * @param {Function} fetchListHandler - The asynchronous function responsible for fetching the list data.
 * @returns {Object} An object containing:
 *   - `data` {Ref<Array>} - Reactive reference to the fetched data array.
 *   - `meta` {Ref<Object>} - Reactive reference to metadata associated with the fetched data.
 *   - `pending` {Ref<boolean>} - Reactive reference indicating the pending state of the fetch operation.
 *   - `fetchData` {Function} - Function to trigger the data fetch operation with the current query.
 *   - `query` {Ref<Object>} - Reactive reference to the current query state.
 *   - `updateQuery` {Function} - Function to update the query state.
 *   - `resetQuery` {Function} - Function to reset the query to its initial state.
 *   - `options` {ComputedRef<Object>} - Computed reference to the current pagination and sorting options.
 *   - `updateOptions` {Function} - Function to update the pagination and sorting options.
 */
export const useFetchListWithoutStoreWithOptions = (fetchListHandler) => {
  const data = ref([]);
  const meta = ref({});

  const {
    pending,
    query,
    updateQuery,
    resetQuery,
    fetchHandlerWithQuery: fetchList,
  } = usePendingWithLocalQuery({
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
