import { getZIndex } from 'vuetify/lib/util/helpers';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

/**
 * Wait a vuetify animation
 *
 * @return {Promise}
 */
export const waitVuetifyAnimation = () => new Promise(resolve => setTimeout(resolve, VUETIFY_ANIMATION_DELAY));

/**
 * Get the maximum z-index among the specified elements.
 *
 * @param {HTMLElement} base - The base element to compare z-index values.
 * @param {number} [minZIndex = 6] - The minimum z-index value to consider.
 * @param {HTMLElement[]} [exclude = []] - Elements to exclude from z-index comparison.
 * @returns {number} The maximum z-index value.
 */
export const getMaxZIndex = (base, minZIndex = 6, exclude = []) => {
  if (!base) {
    return minZIndex;
  }

  const zis = [minZIndex, getZIndex(base)];

  const activeElements = [
    ...document.getElementsByClassName('v-menu__content--active'),
    ...document.getElementsByClassName('v-dialog__content--active'),
  ];

  for (let index = 0; index < activeElements.length; index += 1) {
    if (!exclude.includes(activeElements[index])) {
      zis.push(getZIndex(activeElements[index]));
    }
  }

  return Math.max(...zis);
};
