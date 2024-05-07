import { computed } from 'vue';

import { useComponentInstance } from './vue';

/**
 * Hook to access modal management methods and properties from the Vue instance.
 *
 * @returns {Object} An object containing modal management methods and properties.
 */
export const useModals = () => {
  const vm = useComponentInstance();

  return vm.$modals;
};

/**
 * Custom hook for managing inner modal states and actions within a Vue component.
 *
 * @param {Object} props - The properties passed to the inner modal, including the modal object.
 * @returns {Object} An object containing the modal's computed properties and a method to close the modal.
 */
export const useInnerModal = (props) => {
  const modals = useModals();

  const modal = computed(() => props.modal);
  const config = computed(() => modal.value.config);

  const close = () => modals.hide(props.modal);

  return {
    close,
    modal,
    config,
  };
};
