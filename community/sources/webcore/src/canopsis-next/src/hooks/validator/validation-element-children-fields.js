import { computed, unref } from 'vue';

import { isParent } from '@/helpers/vue-base';

import { useValidator } from './validator';

const EMPTY_VALIDATOR_FIELDS = [];

/**
 * Computes and returns a list of validator fields for the children of a given parent component.
 * This function relies on the injected `$validator` object to access the validation fields.
 * It filters the validator's fields to include only those belonging to the children of the specified parent component.
 * If the `$validator` is not present or the parent is not specified, it returns an empty array.
 *
 * @param {Object|Ref} parent - The parent component or a ref object containing the parent component.
 * @returns {import('vue').ComputedRef<Array>}
 */
export const useValidationElementChildrenFields = (parent) => {
  const validator = useValidator();

  return computed(() => {
    if (!validator) {
      return EMPTY_VALIDATOR_FIELDS;
    }

    const node = unref(parent);
    const nodeEl = node?.$el || node;

    return validator.fields.items.filter(({ vm, el }) => (
      (nodeEl && nodeEl === el) || nodeEl?.contains?.(el) || isParent(vm, node)
    ));
  });
};
