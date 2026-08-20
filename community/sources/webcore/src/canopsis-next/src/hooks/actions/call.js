import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

/**
 * Provides a hook that encapsulates calling an action with success and error handling using popups.
 * This hook uses the `useI18n` for internationalization to fetch localized strings for messages,
 * and `usePopups` for displaying success or error messages in popup format.
 *
 * @param {boolean} throwOnError - Whether to throw the error if it occurs.
 * @param {string} [successPopupType='success'] - The type of popup to display for success messages.
 * @returns {Object} An object containing the `callActionWithPopup` method.
 */
export const useCallActionWithPopup = (throwOnError = false, successPopupType = 'success') => {
  const { t } = useI18n();
  const popups = usePopups();

  const callActionWithPopup = async (action, afterAction, successText, errorText) => {
    const callSuccessText = successText || t('success.default');
    const callErrorText = errorText || t('errors.default');

    try {
      await action();

      popups[successPopupType]({ text: callSuccessText });

      return afterAction?.();
    } catch (err) {
      console.error(err);

      if (throwOnError) {
        throw err;
      }

      return popups.error({ text: callErrorText });
    }
  };

  return { callActionWithPopup };
};
