import { computed } from 'vue';

import { useValidationChildren } from '@/hooks/validator/validation-children';

/**
 * Hook for handling validation header functionalities.
 *
 * This function retrieves the `hasChildrenError` value from the `useValidationChildren` hook and computes
 * the `hasAnyError` boolean value based on it. It also calculates the CSS class for the validation header
 * based on whether there are any errors present.
 *
 * @returns {Object} An object containing:
 * - `hasAnyError`: A computed boolean indicating if there are any errors present.
 * - `validationHeaderClass`: A computed object representing the CSS classes for the validation header.
 */
export const useValidationHeader = () => {
  const { hasChildrenError } = useValidationChildren();

  const hasAnyError = computed(() => hasChildrenError.value);
  const validationHeaderClass = computed(() => ({ 'validation-header': true, 'error--text': this.hasAnyError }));

  return {
    hasAnyError,
    validationHeaderClass,
  };
};
