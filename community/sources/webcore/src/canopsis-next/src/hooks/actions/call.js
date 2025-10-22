import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

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
    const successText = t('success.default');
    const errorText = t('errors.default');

    try {
      await action();

      popups.success({ text: successText });

      return afterAction?.();
    } catch (err) {
      console.error(err);

      return popups.error({ text: errorText });
    }
  };

  return { callActionWithPopup };
};
