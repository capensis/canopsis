import { isObject } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { promisedWait } from '@/helpers/async';

/**
 * Wait a vuetify animation.
 *
 * @return {Promise}
 */
export const waitVuetifyAnimation = () => promisedWait(VUETIFY_ANIMATION_DELAY);

const getElementZIndex = element => Number(window.getComputedStyle(element).zIndex) || 0;

/**
 * Get the maximum z-index among active overlay elements.
 *
 * @param {HTMLElement} base - The base element to compare z-index values.
 * @param {number} [minZIndex = 6] - The minimum z-index value to consider.
 * @param {HTMLElement[]} [exclude = []] - Elements to exclude from z-index comparison.
 * @returns {number} The maximum z-index value.
 */
export const getMaxZIndex = (base, minZIndex = 6, exclude = []) => {
  const zis = [minZIndex];

  if (base) {
    zis.push(getElementZIndex(base));
  }

  const activeElements = [
    ...document.getElementsByClassName('v-overlay--active'),
    ...document.getElementsByClassName('v-menu__content--active'),
    ...document.getElementsByClassName('menuable__content__active'),
    ...document.getElementsByClassName('v-dialog__content--active'),
    ...document.getElementsByClassName('v-navigation-drawer--open'),
  ];

  activeElements.forEach((element) => {
    if (!exclude.includes(element)) {
      zis.push(getElementZIndex(element));
    }
  });

  return Math.max(...zis) + 1;
};

/**
 * Checks whether an element is fully visible across all levels of nested Vuetify tabs.
 *
 * @param {HTMLElement} el - The element to check.
 * @param {HTMLElement} rootEl - The boundary element to stop traversal at.
 * @returns {boolean}
 */
export const isElementVisibleInTabs = (el, rootEl) => {
  let node = el.parentElement;

  while (node && node !== rootEl) {
    if (
      node.classList.contains('v-window-item') &&
      !node.classList.contains('v-window-item--active')
    ) {
      return false;
    }

    node = node.parentElement;
  }

  return true;
};

const getPropertyFromItem = (item, property, fallback) => {
  if (!isObject(item)) {
    return fallback;
  }

  return property.split('.').reduce((acc, key) => acc?.[key], item) ?? fallback;
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
