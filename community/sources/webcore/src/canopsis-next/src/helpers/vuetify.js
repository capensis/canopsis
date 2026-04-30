import { isObject } from 'lodash';
import { getZIndex, getPropertyFromItem } from 'vuetify/lib/util/helpers';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { promisedWait } from '@/helpers/async';

/**
 * Wait a vuetify animation
 *
 * @return {Promise}
 */
export const waitVuetifyAnimation = () => promisedWait(VUETIFY_ANIMATION_DELAY);

/**
 * Get the maximum z-index among the specified elements.
 *
 * @param {HTMLElement} base - The base element to compare z-index values.
 * @param {number} [minZIndex = 6] - The minimum z-index value to consider.
 * @param {HTMLElement[]} [exclude = []] - Elements to exclude from z-index comparison.
 * @returns {number} The maximum z-index value.
 */
export const getMaxZIndex = (base, minZIndex = 6, exclude = []) => {
  const zis = [minZIndex];

  if (base) {
    zis.push(getZIndex(base));
  }

  const activeElements = [
    ...document.getElementsByClassName('v-menu__content--active'),
    ...document.getElementsByClassName('menuable__content__active'),
    ...document.getElementsByClassName('v-dialog__content--active'),
    ...document.getElementsByClassName('v-navigation-drawer--open'),
  ];

  for (let index = 0; index < activeElements.length; index += 1) {
    if (!exclude.includes(activeElements[index])) {
      zis.push(getZIndex(activeElements[index]));
    }
  }

  return Math.max(...zis) + 1;
};

export const getSelectionText = (items, selectionItem, itemValue = 'value', itemText = 'text') => {
  const text = getPropertyFromItem(selectionItem, itemText, selectionItem);

  if (isObject(selectionItem)) {
    return text;
  }

  const value = getPropertyFromItem(selectionItem, itemValue, text);
  const selectionItemObject = items.find(item => getPropertyFromItem(item, itemValue) === value);

  return selectionItemObject ? getPropertyFromItem(selectionItemObject, itemText, text) : text;
};
