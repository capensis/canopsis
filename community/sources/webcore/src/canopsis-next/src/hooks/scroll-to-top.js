import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue';
import { debounce } from 'lodash';

/**
 * Hook for managing scroll-to-top functionality with fullscreen support
 *
 * @param {string} [fulscreenSelector='.view-fullscreen'] - CSS selector for the fullscreen element
 * @returns {Object} Hook interface
 * @returns {import('vue').Ref<boolean>} returns.pageScrolled - Reactive boolean indicating if page is scrolled
 * @returns {Function} returns.scrollToTop - Function to scroll to top with smooth behavior
 */
export const useScrollToTop = (fulscreenSelector = '.view-fullscreen') => {
  const pageScrolled = ref(false);
  let fullscreenElement = null;
  let debouncedCheckScrollPosition = null;

  /**
   * Checks the current scroll position and updates the pageScrolled ref
   * Uses fullscreen element scroll position if available, otherwise uses window scroll position
   */
  const checkScrollPosition = () => {
    pageScrolled.value = fullscreenElement?.scrollTop ?? window.scrollY > 0;
  };

  /**
   * Adds scroll event listeners to the fullscreen element
   * Also triggers initial scroll position check
   */
  const addFullscreenTargetScrollListeners = () => {
    if (fullscreenElement) {
      fullscreenElement.addEventListener('scroll', debouncedCheckScrollPosition);
      debouncedCheckScrollPosition();
    }
  };

  /**
   * Removes scroll event listeners from the fullscreen element
   * Clears the fullscreen element reference and triggers scroll position check
   */
  const removeFullscreenTargetScrollListeners = () => {
    if (fullscreenElement) {
      fullscreenElement.removeEventListener('scroll', debouncedCheckScrollPosition);
      fullscreenElement = null;

      debouncedCheckScrollPosition();
    }
  };

  /**
   * Handles fullscreen change events by updating scroll listeners
   * Removes existing listeners and adds new ones if fullscreen element exists
   */
  const fullscreenListener = () => {
    removeFullscreenTargetScrollListeners();

    nextTick(() => {
      fullscreenElement = document.querySelector(fulscreenSelector);

      if (fullscreenElement) {
        addFullscreenTargetScrollListeners();
      }
    });
  };

  /**
   * Adds event listener for fullscreen change events
   */
  const addFullscreenListener = () => {
    window.addEventListener('fullscreenchange', fullscreenListener);
  };

  /**
   * Removes event listener for fullscreen change events
   */
  const removeFullscreenListener = () => {
    window.removeEventListener('fullscreenchange', fullscreenListener);
  };

  /**
   * Adds scroll event listeners to the document
   */
  const addDocumentScrollListeners = () => {
    document.addEventListener('scroll', debouncedCheckScrollPosition);
  };

  /**
   * Removes scroll event listeners from the document
   */
  const removeDocumentScrollListeners = () => {
    document.removeEventListener('scroll', debouncedCheckScrollPosition);
  };

  /**
   * Scrolls to the top of the page or fullscreen element with smooth behavior
   * Uses fullscreen element if available, otherwise scrolls the window
   */
  const scrollToTop = () => {
    (fullscreenElement ?? window).scrollTo({
      top: 0,
      behavior: 'smooth',
    });
  };

  // Initialize debounced function
  debouncedCheckScrollPosition = debounce(checkScrollPosition, 100);

  // Lifecycle hooks
  onMounted(() => {
    addFullscreenListener();
    addDocumentScrollListeners();
  });

  onBeforeUnmount(() => {
    removeFullscreenListener();
    removeDocumentScrollListeners();
  });

  return {
    pageScrolled,
    scrollToTop,
  };
};
