import { computed, provide } from 'vue';

import Observer from '@/services/observer';

import { promisedTimeout } from '@/helpers/async';

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
 * @param {Function} [options.errorsToValidation = v => v] - The function to convert errors to validation errors.
 * @param {boolean} [options.withTimeout = false] - The property for timeout enabling.
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
export const useSubmittableForm = ({ form, method, scope = null, errorsToValidation = v => v, withTimeout = true }) => {
  const popups = usePopups();
  const { validator, setFormErrors } = useValidationFormErrors(form);
  const { t } = useI18n();

  const afterSubmitObserver = new Observer();

  provide('$afterSubmitObserver', afterSubmitObserver);

  const submitHandler = async (...args) => {
    try {
      const isFormValid = await validator.validateAll();

      if (isFormValid) {
        const data = await method(...args);

        await afterSubmitObserver.notify(data);
      }
    } catch (err) {
      const wasSet = setFormErrors(errorsToValidation(err));

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
      ? (...args) => promisedTimeout(() => submitHandler(...args), 0)
      : submitHandler,
  );

  /**
   * We write custom any errors flag instead of errors.any
   * because we need to keep logic with filtering nullable scope
   */
  const hasAnyErrors = computed(() => !!validator.errors.items.filter(item => (
    item?.scope === scope && validator.errors.vmId === item?.vmId
  )).length);

  const isDisabled = computed(() => {
    if (!validator?.errors) {
      return submitting.value;
    }

    return submitting.value || hasAnyErrors.value;
  });

  return {
    submitting,
    isDisabled,
    submit,
  };
};
