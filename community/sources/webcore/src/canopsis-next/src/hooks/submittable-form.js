import { computed } from 'vue';

import { usePendingHandler } from './query/pending';
import { useValidationFormErrors } from './validator/validation-form-errors';
import { useI18n } from './i18n';
import { usePopups } from './popups';

/**
 * Creates a submittable form handler with validation and error handling.
 *
 * This function integrates form validation, submission handling, and error management. It uses a validator to check
 * the form's validity before submitting. If the form is valid, it proceeds with the submission method provided.
 * Errors during submission are handled gracefully, displaying error messages using a popup system or logging them
 * to the console if they cannot be associated with form fields.
 *
 * @param {Object} options - Configuration options for the submittable form.
 * @param {Object} options.form - The form data object that will be validated.
 * @param {Function} options.method - The submission method to be called if the form is valid.
 * @param {boolean} options.withTimeout - The property for timeout enabling.
 * @returns {Object} An object containing methods and properties to manage the form submission.
 * @example
 * const form = reactive({ username: '', password: '' });
 * const submitMethod = async () => { console.log('Form submitted!'); };
 * const { submit, submitting, isDisabled } = useSubmittableForm({ form, method: submitMethod });
 *
 * // In a Vue component template:
 * <template>
 *   <form @submit.prevent="submit">
 *     <input v-model="form.username" type="text" placeholder="Username">
 *     <input v-model="form.password" type="password" placeholder="Password">
 *     <button :disabled="isDisabled">Submit</button>
 *   </form>
 * </template>
 */
export const useSubmittableForm = ({ form, method, withTimeout = true }) => {
  const popups = usePopups();
  const { validator, setFormErrors } = useValidationFormErrors(form);
  const { t } = useI18n();

  const submitHandler = async (...args) => {
    try {
      const isFormValid = await validator.validateAll();

      if (isFormValid) {
        await method(...args);
      }
    } catch (err) {
      const wasSet = setFormErrors(err);

      if (!wasSet) {
        console.error(err);

        const message = Object.values(err).join('\n');

        popups.error({ text: message || err.details || t('errors.default') });
      }
    }
  };

  const {
    pending: submitting,
    handler: submit,
  } = usePendingHandler(
    /**
     * If `withTimeout` is true, a timeout is set to call `submitHandler` with the provided arguments after 0 ms
     * to avoid combobox lag. Otherwise, `submitHandler` is called directly.
     */
    withTimeout
      ? (...args) => setTimeout(() => submitHandler(...args), 0)
      : submitHandler,
  );

  const isDisabled = computed(() => submitting.value || validator.errors?.any?.());

  return {
    submit,
    submitting,
    isDisabled,
  };
};
