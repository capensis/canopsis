import { computed, unref } from 'vue';

import { useFormContext } from '@/hooks/form';

import { useValidationElementChildrenFields } from './validation-element-children-fields';

/**
 * Checks whether a registered form field currently shows a visible validation error.
 *
 * Checks both `errorMessages` from the Vue component instance (including `exposed` from
 * `setup`) and vee-validate's `errors` array. Some components expose an empty external
 * `errorMessages` prop while rendering errors registered internally by vee-validate, so
 * neither source may mask the other.
 *
 * @param {Object} [field={}] - Validator field state with its bound component instance.
 * @param {string[]} [field.errors] - Validation error messages from the form context.
 * @param {Object} [field.instance] - Vue component instance of the field.
 * @returns {boolean} `true` if the field has at least one visible error message.
 */
export const hasVisibleFieldError = ({ errors, instance } = {}) => {
  const errorMessages = unref(instance?.errorMessages ?? instance?.$?.exposed?.errorMessages);

  return (Array.isArray(errorMessages) && errorMessages.length > 0) || errors?.length > 0;
};

/**
 * Provides validation functionalities for the children of a given parent component.
 *
 * This function utilizes the injected `$validator` object to perform validations on the children fields of
 * the specified parent component. It computes a boolean value indicating whether any of the children fields
 * have validation errors and provides a method to validate all children fields.
 *
 * @param {Object|Ref} parent - The parent component or a ref object containing the parent component.
 * @returns {Object} An object containing:
 * - `childrenFields`: A computed reference to the list of validator fields for the children of the parent component.
 * - `hasChildrenError`: A computed boolean indicating if any of the children fields have validation errors.
 * - `validateChildren`: A function that triggers validation for all children fields. Returns a Promise if the
 * validator is present, otherwise `undefined`.
 */
export const useValidationElementChildren = (parent) => {
  const formContext = useFormContext();

  const childrenFields = useValidationElementChildrenFields(parent);

  const hasChildrenError = computed(() => childrenFields.value.some(hasVisibleFieldError));

  const validateChildren = () =>
    childrenFields.value.forEach(({ path } = {}) => formContext.validateField(path));

  return {
    childrenFields,
    hasChildrenError,
    validateChildren,
  };
};
