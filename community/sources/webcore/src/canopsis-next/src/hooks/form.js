import { computed } from 'vue';

import { useComponentModel } from './vue';

/**
 * Hook for update model value
 *
 * @param {Object} props
 * @param {Function} emit
 * @return {import('vue').WritableComputedRef}
 */
export const useModelValue = (props, emit) => {
  const { prop, event } = useComponentModel();

  return computed({
    set(value) {
      emit(event, value);
    },
    get() {
      return props[prop];
    },
  });
};
