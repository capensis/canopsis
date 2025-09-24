import { ref, unref } from 'vue';
import { cloneDeep, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import { useModals } from './modals';

/**
 * Hook for confirmable form functionality that shows a confirmation modal
 * when attempting to perform an action on a changed form.
 *
 * @param {Object} options - Configuration options for the hook
 * @param {Function} options.method - The method to wrap with confirmation logic
 * @param {string} [options.modalName = MODALS.confirmation] - The modal name to show
 * @param {Function} [options.comparator = isEqual] - Function to compare form values
 * @param {boolean} [options.cloning = false] - Whether to clone the initial form value
 * @param {Object} [options.form] - The form object to track (when cloning is true)
 * @returns {Object} An object containing the confirmation method
 *
 * @example
 * // In a Vue component setup function
 * const removeItem = () => {
 *   // original remove logic
 * };
 *
 * const { confirmAction } = useConfirmableForm({
 *   method: removeItem,
 *   cloning: true,
 *   form: props.form,
 * });
 *
 * // Use confirmAction instead of removeItem directly
 * const handleRemove = () => {
 *   confirmAction();
 * };
 */
export const useConfirmableForm = ({
  method,
  modalName = MODALS.confirmation,
  comparator = isEqual,
  cloning = false,
  form,
} = {}) => {
  const modals = useModals();
  const confirmationModalId = uid('modal');

  const originalForm = ref(cloning && form ? cloneDeep(unref(form)) : null);

  /**
   * Confirms an action by checking if the form has changed.
   * If changed, shows a confirmation modal; otherwise executes the action directly.
   *
   * @param {Function} [action] - The action to execute (defaults to the configured method)
   * @param {*} [currentForm] - The current form value to compare (when cloning is true)
   */
  const confirmAction = (action = method, currentForm = form) => {
    let equal = true;

    if (cloning && originalForm.value && currentForm) {
      equal = comparator(unref(currentForm), originalForm.value);
    }

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
    } else {
      action?.();
    }
  };

  /**
   * Updates the original form value (useful when form is reset or changed externally)
   *
   * @param {*} newForm - The new form value to store as original
   */
  const updateOriginalForm = (newForm) => {
    if (cloning) {
      originalForm.value = cloneDeep(unref(newForm));
    }
  };

  return {
    confirmAction,
    updateOriginalForm,
    originalForm: originalForm.value,
  };
};
