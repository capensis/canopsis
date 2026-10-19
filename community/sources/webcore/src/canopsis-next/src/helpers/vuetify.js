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

/**
 * Scrolls to the first invalid field visible in the current tab.
 *
 * @param {HTMLElement} [rootEl] Form or modal root element.
 */
export const scrollToFirstFormError = (rootEl) => {
  if (!rootEl?.querySelector) {
    return;
  }

  const errorElements = [...rootEl.querySelectorAll('.v-input--error, .v-field--error')];

  if (!errorElements.length) {
    return;
  }

  const visibleError = errorElements.find(el => isElementVisibleInTabs(el, rootEl));

  if (visibleError) {
    visibleError.scrollIntoView({ behavior: 'smooth', block: 'center' });

    return;
  }

  rootEl.querySelector('.v-card-text, .v-card__text')?.scrollTo({ top: 0, behavior: 'smooth' });
};

/**
 * Reads a possibly nested property from an item using a dotted path.
 *
 * @param {*} item Source value.
 * @param {string} property Property path, for example `text` or `user.name`.
 * @param {*} fallback Value returned when `item` is not an object or the path is missing.
 * @returns {*} Property value or `fallback`.
 */
const getPropertyFromItem = (item, property, fallback) => {
  if (!isObject(item)) {
    return fallback;
  }

  return property.split('.').reduce((acc, key) => acc?.[key], item) ?? fallback;
};

/**
 * Resolves the display text for a selected Vuetify select/combobox value.
 *
 * If `selectionItem` is an object, its `itemText` property is returned.
 * If it is a primitive, the matching item is found in `items` by `itemValue`
 * and that item's `itemText` is returned.
 *
 * @param {Array} items Select items.
 * @param {*} selectionItem Selected object or primitive value.
 * @param {string} [itemValue='value'] Value property path on each item.
 * @param {string} [itemText='text'] Text property path on each item.
 * @returns {*} Display text for the selection.
 */
export const getSelectionText = (items, selectionItem, itemValue = 'value', itemText = 'text') => {
  const text = getPropertyFromItem(selectionItem, itemText, selectionItem);

  if (isObject(selectionItem)) {
    return text;
  }

  const value = getPropertyFromItem(selectionItem, itemValue, text);
  const selectionItemObject = items.find(item => getPropertyFromItem(item, itemValue) === value);

  return selectionItemObject ? getPropertyFromItem(selectionItemObject, itemText, text) : text;
};
