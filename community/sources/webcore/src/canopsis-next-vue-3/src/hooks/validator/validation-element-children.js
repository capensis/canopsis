import { computed } from 'vue';

import { useFormErrors, useFormContext } from '@/hooks/form';

import { useValidationElementChildrenFields } from './validation-element-children-fields';

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
  const formErrors = useFormErrors();
  const formContext = useFormContext();
  const childrenFields = useValidationElementChildrenFields(parent);

  const hasChildrenError = computed(() => {
    if (!Object.values(formErrors.value).length) {
      return false;
    }

    return childrenFields.value.some(({ errors } = {}) => errors?.length > 0);
  });

  const validateChildren = () =>
    childrenFields.value.forEach(({ path } = {}) => formContext?.validateField?.(path));

  return {
    childrenFields,
    hasChildrenError,
    validateChildren,
  };
};
