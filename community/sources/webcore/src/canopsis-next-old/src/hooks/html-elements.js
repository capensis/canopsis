import { computed, unref } from 'vue';

/**
 * Function to retrieve a specific element within a table element.
 *
 * @param {Object} options - The options object.
 * @param {Object} options.parentElement - The table element to search within.
 * @param {string} options.selector - The CSS selector to find the element.
 * @returns {Object} - The found element.
 */
export const useHTMLElement = ({ parentElement, selector }) => {
  const element = computed(() => {
    const unwrappedParentElement = unref(parentElement);

    if (!unwrappedParentElement) {
      return null;
    }

    return (unwrappedParentElement.$el ?? unwrappedParentElement).querySelector(unref(selector));
  });

  return {
    element,
  };
};
/**
 * Function to retrieve a specific elements within a table element.
 *
 * @param {Object} options - The options object.
 * @param {Object} options.parentElement - The parent element to search within.
 * @param {string} options.selector - The CSS selector to find the element.
 * @returns {Object} - The found elements.
 */
export const useHTMLElements = ({ parentElement, selector }) => {
  const elements = computed(() => {
    const unwrappedParentElement = unref(parentElement);

    if (!unwrappedParentElement) {
      return [];
    }

    return (unwrappedParentElement.$el ?? unwrappedParentElement).querySelectorAll(unref(selector));
  });

  return {
    elements,
  };
};
