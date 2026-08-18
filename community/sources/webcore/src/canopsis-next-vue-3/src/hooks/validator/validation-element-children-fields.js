import { computed, unref } from 'vue';

import { isParent } from '@/helpers/vue-base';

import { useFormContext } from '@/hooks/form';

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
  const formContext = useFormContext();

  return computed(() => {
    if (!formContext) {
      return EMPTY_VALIDATOR_FIELDS;
    }

    const fieldsWithInstances = formContext.getAllPathStatesWithInstances();

    const node = unref(parent);
    const nodeEl = node?.$el;

    return fieldsWithInstances.filter(({ instance } = {}) =>
      (nodeEl && nodeEl === instance?.$el) || nodeEl?.contains?.(instance?.$el) || isParent(instance, node),
    );
  });
};
