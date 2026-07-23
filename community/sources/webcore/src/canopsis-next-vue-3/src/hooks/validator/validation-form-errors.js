import { has } from 'lodash';
import { unref } from 'vue';

import { useValidator } from './validator';

/**
 * Provides utilities for managing form validation errors using a validator object.
 *
 * @param {Object|import('vue').Ref<Object>} form - The form object which contains the fields to validate.
 * @returns {Object} An object containing the validator instance and a method to set form errors.
 */
export const useValidationFormErrors = (form) => {
  const validator = useValidator();

  /**
   * Get errors for existing fields in current form.
   *
   * @param {Object} errors
   * @returns {[string, string][]}
   */
  const getExistingFieldsErrors = (errors = {}) => {
    const unwrappedForm = unref(form);

    return Object.entries(errors).filter(
      ([field]) => has(unwrappedForm, field) || Array.isArray(unwrappedForm),
    );
  };

  /**
   * Add existing field errors to validator errors.
   *
   * @param {[string, string][]} existingFieldsErrors
   * @param {string} [scope]
   */
  const addExistingFieldsErrors = (existingFieldsErrors, scope) => {
    validator.errors.add(
      existingFieldsErrors.map(([field, msg]) => ({
        field: scope ? `${scope}.${field}` : field,
        msg,
      })),
    );
  };

  /**
   * Set form errors from response error and returns true if form errors exist.
   *
   * @param {Object} [errors={}]
   * @param {string} [scope]
   * @returns {boolean}
   */
  const setFormErrors = (errors = {}, scope) => {
    if (!validator) {
      return false;
    }

    const existingFieldErrors = getExistingFieldsErrors(errors);

    if (existingFieldErrors.length) {
      addExistingFieldsErrors(existingFieldErrors, scope);

      return true;
    }

    return false;
  };

  return {
    validator,
    setFormErrors,
  };
};
