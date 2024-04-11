import { useComponentInstance } from '../vue';

import { useValidationElementChildren } from './useValidationElementChildren';

/**
 * Hook for validating the children of the current component instance.
 *
 * @returns {Object}
 */
export const useValidationChildren = () => {
  const instance = useComponentInstance();

  return useValidationElementChildren(instance);
};
