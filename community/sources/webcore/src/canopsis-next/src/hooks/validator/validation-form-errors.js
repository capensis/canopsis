import { computed, unref } from 'vue';
import { has, keyBy } from 'lodash';

import { useValidator } from './validator';

/**
 * Provides utilities for managing form validation errors using a validator object.
 * This function encapsulates methods to set and manage errors based on a given form structure.
 *
 * @param {Object} form - The form object which contains the fields to validate.
 * @returns {Object} An object containing the validator instance and a method to set form errors.
 *
 * @example
 * // Assuming `form` is an object that represents your form fields
 * const form = {
 *   username: '',
 *   password: ''
 * };
 *
 * // Using the useValidationFormErrors function
 * const { validator, setFormErrors } = useValidationFormErrors(form);
 *
 * // Example error object that might be returned from a server on form submission
 * const errorsFromServer = {
 *   username: 'Username is required',
 *   password: 'Password must be at least 6 characters'
 * };
 *
 * // Setting form errors and checking if there are any
 * const hasErrors = setFormErrors(errorsFromServer);
 * if (hasErrors) {
 *   console.log('There are errors in the form');
 * }
 */
export const useValidationFormErrors = (form) => {
  const validator = useValidator();

  const fieldsByName = computed(() => keyBy(validator.fields.items, 'name'));

  /**
   * Get errors for exists fields in current form
   *
   * @param {any} errors
   * @returns {[string, string][]}
   */
  const getExistsFieldsErrors = errors => Object.entries(errors)
    .filter(([field]) => fieldsByName.value[field] || has(unref(form), field));

  /**
   * Add exists fields errors to validator errors
   *
   * @param {[string, string][]} existsFieldsErrors
   */
  const addExistsFieldsErrors = (existsFieldsErrors) => {
    validator.errors.add(existsFieldsErrors.map(([field, msg]) => ({ field, msg })));
  };

  /**
   * Set form errors from response error and returns true if form errors exists
   *
   * @param {any} [errors = {}]
   * @return {boolean}
   */
  const setFormErrors = (errors = {}) => {
    if (!validator) {
      return false;
    }

    const existFieldErrors = getExistsFieldsErrors(errors);

    if (existFieldErrors.length) {
      addExistsFieldsErrors(existFieldErrors);

      return true;
    }

    return false;
  };

  return {
    validator,
    setFormErrors,
  };
};
