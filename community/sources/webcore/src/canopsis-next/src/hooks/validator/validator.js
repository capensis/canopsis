import { inject } from 'vue';

import { useComponentInstance } from '@/hooks/vue';

/**
 * Injects the validator object into the component.
 *
 * This function uses Vue's `inject` method to retrieve the validator object, identified by the key `$validator`,
 * from an ancestor component. The validator object is typically provided by a vee-validate.
 *
 * @returns {Object}
 */
export const useValidator = () => {
  const vm = useComponentInstance();

  return vm.$validator || inject('$validator');
};
