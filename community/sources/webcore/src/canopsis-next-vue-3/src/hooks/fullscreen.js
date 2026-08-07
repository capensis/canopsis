import { ref } from 'vue';

const isFullscreen = ref(false);
let cleanupClass = null;
let fullscreenCallback = null;

const updateFullscreenState = () => {
  const value = Boolean(document.fullscreenElement);

  isFullscreen.value = value;
  fullscreenCallback?.(value);

  if (!value && cleanupClass) {
    cleanupClass();
    cleanupClass = null;
    fullscreenCallback = null;
  }
};

document.addEventListener('fullscreenchange', updateFullscreenState);

/**
 * Hook that provides access to fullscreen functionality.
 *
 * @returns {Object}
 */
export const useFullscreen = () => {
  const request = async (element, { fullscreenClass, callback } = {}) => {
    if (!element) {
      return;
    }

    fullscreenCallback = callback;

    if (fullscreenClass) {
      element.classList.add(fullscreenClass);
      cleanupClass = () => element.classList.remove(fullscreenClass);
    }

    await element.requestFullscreen?.();
    updateFullscreenState();
  };

  const exit = async () => {
    if (!document.fullscreenElement) {
      return;
    }

    await document.exitFullscreen?.();
    updateFullscreenState();
  };

  const toggle = (element, options) => (isFullscreen.value ? exit() : request(element, options));

  return {
    isFullscreen,
    request,
    exit,
    toggle,
  };
};
