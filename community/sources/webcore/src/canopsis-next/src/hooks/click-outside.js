import { inject, onBeforeUnmount } from 'vue';

/**
 * Injects the click outside functionality from Vue's injection system.
 * This function is a composition API utility that provides access to a global click outside directive.
 *
 * @returns {ClickOutside} Returns the injected click outside handler.
 */
export const useInjectClickOutside = () => inject('$clickOutside');

/**
 * Registers a click outside handler and ensures it is unregistered when the component is unmounted.
 * This function utilizes the injected click outside functionality to manage the lifecycle
 * of the click outside event handler.
 *
 * @param {Function} handler - The function to be executed when a click outside event is detected.
 */
export const useRegisterClickOutsideHandler = (handler) => {
  const clickOutside = useInjectClickOutside();
  clickOutside.register(handler);

  onBeforeUnmount(() => {
    clickOutside.unregister(handler);
  });
};
