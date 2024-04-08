import { computed, inject, unref } from 'vue';

import { isParent } from '@/helpers/vue-base';

import { useComponentInstance } from '../vue';

const EMPTY_VALIDATOR_FIELDS = [];

/**
 * Injects the validator object into the component.
 *
 * This function uses Vue's `inject` method to retrieve the validator object, identified by the key `$validator`,
 * from an ancestor component. The validator object is typically provided by a vee-validate.
 *
 * @returns {Object}
 */
export const useInjectValidator = () => inject('$validator');

/**
 * Computes and returns a list of validator fields for the children of a given parent component.
 * This function relies on the injected `$validator` object to access the validation fields.
 * It filters the validator's fields to include only those belonging to the children of the specified parent component.
 * If the `$validator` is not present or the parent is not specified, it returns an empty array.
 *
 * @param {Object|Ref} parent - The parent component or a ref object containing the parent component.
 * @returns {import('vue').ComputedRef<Array>
 */
export const useElementChildrenFields = (parent) => {
  const validator = useInjectValidator();

  return computed(() => {
    if (!validator) {
      return EMPTY_VALIDATOR_FIELDS;
    }

    const node = unref(parent);

    return validator.fields.items.filter(({ vm }) => isParent(vm, node));
  });
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
export const useElementChildrenValidation = (parent) => {
  const validator = useInjectValidator();
  const childrenFields = useElementChildrenFields(parent);

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

/**
 * Hook for validating the children of the current component instance.
 *
 * @returns {Object}
 */
export const useValidationChildren = () => {
  const instance = useComponentInstance();

  return useElementChildrenValidation(instance);
};
