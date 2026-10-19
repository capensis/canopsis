import { useComponentInstance } from './vue';

/**
 * Hook that provides access to fullscreen functionality.
 *
 * @returns {Object} The fullscreen service object with methods to control fullscreen mode
 * @returns {Function} returns.enter - Method to enter fullscreen mode
 * @returns {Function} returns.exit - Method to exit fullscreen mode
 * @returns {Function} returns.toggle - Method to toggle fullscreen mode
 * @returns {boolean} returns.isFullscreen - Current fullscreen state
 */
export const useFullscreen = () => {
  const instance = useComponentInstance();

  return instance.$fullscreen;
};
