import { computed } from 'vue';

import { useValidationChildren } from '@/hooks/validator/validation-children';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

/**
 * Hook for handling validation header functionalities.
 *
 * @param {Object|import('vue').Ref} [parent] Optional content ref for scoped child validation.
 * @returns {Object} Validation state and css classes.
 */
export const useValidationHeader = (parent) => {
  const validationChildren = parent
    ? useValidationElementChildren(parent)
    : useValidationChildren();
  const { hasChildrenError } = validationChildren;

  const hasAnyError = computed(() => hasChildrenError.value);
  const validationHeaderClass = computed(() => ({
    'validation-header': true,
    'error--text': hasAnyError.value,
    'text-error': hasAnyError.value,
  }));

  return {
    hasAnyError,
    validationHeaderClass,
  };
};
