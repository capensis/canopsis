import { computed } from 'vue';

import { useValidationElementChildrenFields } from './useValidationElementChildrenFields';
import { useInjectValidator } from './useInjectValidator';

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
  const validator = useInjectValidator();
  const childrenFields = useValidationElementChildrenFields(parent);

  const hasChildrenError = computed(() => {
    if (validator && validator.errors.any()) {
      return childrenFields.value.some(field => validator.errors.has(field.name));
    }

    return false;
  });

  const validateChildren = () => validator?.validateAll?.(childrenFields.value);

  return {
    childrenFields,
    hasChildrenError,
    validateChildren,
  };
};
