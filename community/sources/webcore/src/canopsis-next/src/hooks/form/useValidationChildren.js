import { computed, inject, isRef } from 'vue';

import { isParent } from '@/helpers/vue-base';

import { useComponentInstance } from '../vue';

const EMPTY_VALIDATOR_FIELDS = [];

export const useInjectValidator = () => inject('$validator');

export const useElementChildrenFields = (parent) => {
  const validator = useInjectValidator();

  return computed(() => {
    if (!validator) {
      return EMPTY_VALIDATOR_FIELDS;
    }

    const node = isRef(parent) ? parent.value : parent;

    return validator.fields.items.filter(({ vm }) => isParent(vm, node));
  });
};

export const useElementChildrenValidation = (parent) => {
  const validator = useInjectValidator();
  const childrenFields = useElementChildrenFields(parent);

  const hasChildrenError = computed(() => {
    if (validator && validator.errors.any()) {
      return childrenFields.value.some(field => validator.errors.has(field.name));
    }

    return false;
  });

  const validateChildren = options => validator?.validateAll?.(childrenFields.value, options);

  return {
    childrenFields,
    hasChildrenError,
    validateChildren,
  };
};

export const useValidationChildren = () => {
  const instance = useComponentInstance();

  return useElementChildrenValidation(instance);
};
