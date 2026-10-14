import { unref } from 'vue';

import { useFormContext } from '@/hooks/form';
import { getMatchedServerFieldErrors } from '@/hooks/form/server-errors';

/**
 * Provides utilities for mapping server errors onto the current form context.
 *
 * @param {Object|import('vue').Ref<Object>} [form] Form values used to match server fields.
 * @param {Object|import('vue').Ref<Object>} [providedFormContext] Form context used to set field errors.
 * @returns {Object} An object containing a method to set form errors.
 *
 * @example
 * const form = { username: '', password: '' };
 * const { setFormErrors } = useValidationFormErrors(form);
 *
 * const hasErrors = setFormErrors({
 *   username: 'Username is required',
 *   password: 'Password must be at least 6 characters',
 * });
 */
export const useValidationFormErrors = (form, providedFormContext) => {
  const defaultFormContext = useFormContext();

  /**
   * Set form errors from a server response and returns true if any field errors were applied.
   *
   * @param {Object} [errors={}]
   * @param {string} [scope]
   * @returns {boolean}
   */
  const setFormErrors = (errors = {}, scope) => {
    const formContext = unref(providedFormContext) ?? defaultFormContext;
    const { addError, values } = formContext;

    const existingFieldErrors = getMatchedServerFieldErrors(errors, {
      formValues: form ?? values,
      formContext,
    });

    if (!existingFieldErrors.length) {
      return false;
    }

    existingFieldErrors.forEach(([field, msg]) => {
      addError(scope ? `${scope}.${field}` : field, msg);
    });

    return true;
  };

  return {
    setFormErrors,
  };
};
