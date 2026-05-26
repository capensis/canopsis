import { computed, provide, unref } from 'vue';

import Observer from '@/services/observer';

import { promisedTimeout } from '@/helpers/async';
import { isElementVisibleInTabs } from '@/helpers/vuetify';

import { usePendingHandler } from './query/pending';
import { useValidationFormErrors } from './validator/validation-form-errors';
import { useI18n } from './i18n';
import { usePopups } from './popups';
import { useComponentInstance } from './vue';

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
 * const { submit, submitting, submitLabel } = useSubmittableForm({ form, item, method: submitMethod });
 *
 * // In a Vue component template:
 * <template>
 *   <form @submit.prevent="submit">
 *     <input v-model="form.username" type="text" placeholder="Username">
 *     <input v-model="form.password" type="password" placeholder="Password">
 *     <button :disabled="submitting">Submit</button>
 *   </form>
 * </template>
 */
export const useSubmittableForm = ({
  form,
  item,
  method,
  errorsToValidation = v => v,
  withTimeout = true,
}) => {
  const popups = usePopups();
  const { validator, setFormErrors } = useValidationFormErrors(form);
  const { t } = useI18n();
  const instance = useComponentInstance();

  const afterSubmitObserver = new Observer();

  provide('$afterSubmitObserver', afterSubmitObserver);

  /**
   * Scrolls to the first error field visible in the currently active tab.
   * If no errors exist in the active tab, scrolls the modal content area to the top
   * so the user can see the tab indicator highlighting the tab that contains errors.
   */
  const scrollToFirstError = () => {
    const errorFieldNames = new Set(validator.errors.items.map(e => e.field));

    const allErrorElements = validator.fields.items
      .filter(field => errorFieldNames.has(field.name) && field.el)
      .map(field => field.el);

    if (!allErrorElements.length) {
      return;
    }

    const rootEl = instance.$el;

    const activeTabErrors = allErrorElements.filter(el => isElementVisibleInTabs(el, rootEl));

    if (activeTabErrors.length > 0) {
      activeTabErrors[0].scrollIntoView({ behavior: 'smooth', block: 'center' });

      return;
    }

    rootEl?.querySelector('.v-card__text')?.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const submitHandler = async (...args) => {
    try {
      const isFormValid = await validator.validateAll();

      if (!isFormValid) {
        scrollToFirstError();

        return;
      }

      const data = await method(...args);

      await afterSubmitObserver.notify(data);
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

  const isNew = computed(() => !unref(item)?._id);

  const submitLabel = computed(() => (isNew.value ? t('common.create') : t('common.save')));

  return {
    submitting,
    submit,
    isNew,
    submitLabel,
  };
};
