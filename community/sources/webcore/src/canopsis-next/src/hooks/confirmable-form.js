import { cloneDeep, isEqual } from 'lodash';
import { onBeforeMount, unref } from 'vue';

import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import { useModals } from './modals';

/**
 * Hook to manage form confirmation when data has changed.
 * It checks if the form data has changed from its original state and,
 * if so, prompts the user with a confirmation modal.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Ref} options.form - Vue ref object containing the form data.
 * @param {string} [options.modalName=MODALS.confirmation] - Name of the modal to show for confirmation.
 * @param {Function} [options.comparator=isEqual] - Function used to compare the current form state with the original.
 * @returns {Object} An object containing the authToken ref and confirmation method.
 *
 * @example
 * export default {
 *   setup(props, { emit }) {
 *     const { authToken, confirmAction } = useConfirmableForm({
 *       form: toRef(props, 'form'),
 *     });
 *
 *     const removeWebhook = () => {
 *       confirmAction(() => {
 *         emit('remove');
 *       });
 *     };
 *
 *     return {
 *       authToken,
 *       removeWebhook,
 *     };
 *   },
 * };
 */
export const useConfirmableForm = ({
  form,
  action,
  modalName = MODALS.confirmation,
  comparator = isEqual,
  cloning = false,
} = {}) => {
  const modals = useModals();

  let originalForm = {};

  const confirmationModalId = uid('confirmation');

  const confirmAction = () => {
    const equal = comparator(unref(form), originalForm);

    if (!equal) {
      modals.show({
        id: confirmationModalId,
        name: modalName,
        dialogProps: {
          persistent: true,
        },
        config: {
          action,
        },
      });

      return;
    }

    action();
  };

  onBeforeMount(() => {
    if (cloning) {
      originalForm = cloneDeep(unref(form));
    }
  });

  return {
    confirmAction,
  };
};
