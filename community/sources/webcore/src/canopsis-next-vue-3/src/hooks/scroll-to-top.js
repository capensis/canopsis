import { debounce } from 'lodash';
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue';

/**
 * Hook for managing scroll-to-top functionality with fullscreen support.
 *
 * @param {string} [fullscreenSelector='.view-fullscreen'] - CSS selector for the fullscreen element.
 * @returns {Object}
 */
export const useScrollToTop = (fullscreenSelector = '.view-fullscreen') => {
  const pageScrolled = ref(false);
  let fullscreenElement = null;

  const checkScrollPosition = () => {
    pageScrolled.value = fullscreenElement?.scrollTop > 0 || window.scrollY > 0;
  };

  const debouncedCheckScrollPosition = debounce(checkScrollPosition, 100);

  const addFullscreenTargetScrollListeners = () => {
    if (!fullscreenElement) {
      return;
    }

    fullscreenElement.addEventListener('scroll', debouncedCheckScrollPosition);
    debouncedCheckScrollPosition();
  };

  const removeFullscreenTargetScrollListeners = () => {
    if (!fullscreenElement) {
      return;
    }

    fullscreenElement.removeEventListener('scroll', debouncedCheckScrollPosition);
    fullscreenElement = null;
    debouncedCheckScrollPosition();
  };

  const fullscreenListener = () => {
    removeFullscreenTargetScrollListeners();

    nextTick(() => {
      fullscreenElement = document.querySelector(fullscreenSelector);
      addFullscreenTargetScrollListeners();
    });
  };

  const scrollToTop = () => {
    (fullscreenElement ?? window).scrollTo({
      top: 0,
      behavior: 'smooth',
    });
  };

  onMounted(() => {
    window.addEventListener('fullscreenchange', fullscreenListener);
    document.addEventListener('scroll', debouncedCheckScrollPosition);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('fullscreenchange', fullscreenListener);
    document.removeEventListener('scroll', debouncedCheckScrollPosition);
    removeFullscreenTargetScrollListeners();
  });

  return {
    pageScrolled,
    scrollToTop,
  };
};
