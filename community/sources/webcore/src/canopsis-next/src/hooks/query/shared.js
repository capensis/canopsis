import { usePendingHandler } from './pending';
import { useLocalQuery } from './local-query';

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
    handler: fetchWithPending,
  } = usePendingHandler(fetchHandler, initialPending);

  const queryData = useLocalQuery({
    initialQuery,
    comparator,
    onUpdate: fetchWithPending,
  });

  const fetchHandlerWithQuery = (requestQuery = queryData.query.value) => fetchWithPending(requestQuery);

  return {
    ...queryData,

    pending,

    fetchHandlerWithQuery,
  };
};
