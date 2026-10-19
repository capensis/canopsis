import { cloneDeep, isEqual } from 'lodash';
import { provide, ref, unref } from 'vue';

import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import { useRegisterClickOutsideHandler } from './click-outside';
import { useModals } from './modals';

/**
 * Hook to manage modal confirmation on unsaved form changes when attempting to close the modal.
 * It checks if the form data has changed from its original state and,
 * if so, prompts the user with a confirmation modal.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Function} options.submit - Function to call when the user confirms the action.
 * @param {Ref} options.form - Vue ref object containing the form data.
 * @param {string} [options.modalName=MODALS.clickOutsideConfirmation] - Name of the modal to show for confirmation.
 * @param {Function} options.close - Function to call when the modal is closed without changes or after confirmation.
 * @param {Function} [options.comparator=isEqual] - Function used to compare the current form state with the original.
 * @returns {Object}
 *
 * @example
 * const form = ref({ name: '' });
 *
 * useFormConfirmableCloseModal({
 *   form,
 *   submit,
 *   close,
 * });
 */
export const useFormConfirmableCloseModal = ({
  submit,
  form,
  modalName = MODALS.clickOutsideConfirmation,
  close,
  comparator = isEqual,
} = {}) => {
  const modals = useModals();

  const originalForm = ref(cloneDeep(unref(form)));

  const confirmationModalId = uid('modal');

  const clickOutsideHandlerMethod = () => {
    const equal = comparator(unref(form), originalForm.value);

    const hasOpenedModal = modals.isModalOpenedById(confirmationModalId);

    if (!equal && !hasOpenedModal) {
      modals.show({
        id: confirmationModalId,
        name: modalName,
        dialogProps: {
          persistent: true,
        },
        config: {
          action: (confirmed) => {
            if (confirmed) {
              return submit?.();
            }

            return close?.();
          },
        },
      });
    }

    return equal;
  };

  const updateOriginalForm = () => (originalForm.value = cloneDeep(unref(form)));
  const updateOriginalFormField = (field, value) => (originalForm.value[field] = value);

  provide('$closeModal', () => {
    if (clickOutsideHandlerMethod()) {
      close?.();
    }
  });

  useRegisterClickOutsideHandler(clickOutsideHandlerMethod);

  return {
    updateOriginalForm,
    updateOriginalFormField,
  };
};
