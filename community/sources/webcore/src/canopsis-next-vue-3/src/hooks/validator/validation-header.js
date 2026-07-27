import { computed } from 'vue';

import { useValidationChildren } from '@/hooks/validator/validation-children';

/**
 * Hook for handling validation header functionalities.
 *
 * @returns {Object} Validation state and css classes.
 */
export const useValidationHeader = () => {
  const { hasChildrenError } = useValidationChildren();

  const hasAnyError = computed(() => hasChildrenError.value);
  const validationHeaderClass = computed(() => ({
    'validation-header': true,
    'error--text': hasAnyError.value,
  }));

  return {
    hasAnyError,
    validationHeaderClass,
  };
};
