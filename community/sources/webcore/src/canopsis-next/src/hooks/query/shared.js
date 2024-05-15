import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

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

  return {
    ...queryData,
    pending,
  };
};

/**
 * Provides a hook that encapsulates calling an action with success and error handling using popups.
 * This hook uses the `useI18n` for internationalization to fetch localized strings for messages,
 * and `usePopups` for displaying success or error messages in popup format.
 *
 * @returns {Object} An object containing the `callActionWithPopup` method.
 */
export const useCallActionWithPopup = () => {
  const { t } = useI18n();
  const popups = usePopups();

  const callActionWithPopup = async (action, afterAction) => {
    try {
      await action();

      popups.success({ text: t('success.default') });

      return afterAction();
    } catch (err) {
      console.error(err);

      return popups.error({ text: t('errors.default') });
    }
  };

  return { callActionWithPopup };
};
