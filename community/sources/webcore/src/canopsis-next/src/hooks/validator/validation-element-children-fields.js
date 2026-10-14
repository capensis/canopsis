import { computed, unref } from 'vue';

import { isParent } from '@/helpers/vue-base';

import { useFormContext } from '@/hooks/form';

/**
 * Resolves a Vue component instance or a DOM node to a native Element.
 *
 * @param {Object|Element|null|undefined} node - Vue component instance (`$el`) or a DOM node.
 * @returns {Element|null}
 */
const getElement = (node) => {
  const el = node?.$el ?? node;

  return el instanceof Element ? el : null;
};

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
    const fieldsWithInstances = formContext.getAllPathStatesWithInstances();

    const node = unref(parent);
    const nodeEl = getElement(node);

    return fieldsWithInstances.filter(({ instance } = {}) => {
      const instanceEl = getElement(instance);

      return (
        (nodeEl && nodeEl === instanceEl) ||
        nodeEl?.contains?.(instanceEl) ||
        isParent(instance, node)
      );
    });
  });
};
