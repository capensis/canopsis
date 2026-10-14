import { isNil, isObject } from 'lodash';

/**
 * @typedef {Object} SlotItemPayload
 * Item payload from a Vuetify select/autocomplete/combobox slot (`#item`, `#selection`, `#chip`).
 *
 * @property {string|number} [title] Normalized item title from Vuetify's internal item.
 * @property {*} [raw] Original item passed into `items`.
 * @property {*} [name] Legacy display field.
 * @property {*} [text] Legacy display field.
 * @property {*} [value] Item value used as a last-resort title.
 */

/**
 * @typedef {Object} SlotItemProps
 * Slot props passed to field item/selection/chip slots.
 *
 * @property {SlotItemPayload} [internalItem] Vuetify internal item (`title`, `value`, `raw`).
 * @property {SlotItemPayload|*} [item] Slot item: internal item, `{ raw }`, or the original value.
 */

/**
 * Resolves a display title from a Vuetify item/selection/chip slot payload.
 *
 * Prefers `internalItem.title`, then `item.title`, then legacy fields on the raw item
 * (`name`, `text`, `title`, `value`).
 *
 * @param {SlotItemProps} [slotProps={}] Slot payload from Vuetify or a wrapper (`getLegacySlotProps`).
 * @returns {string} Display title, or an empty string when none can be resolved.
 */
export const getSlotItemTitle = (slotProps = {}) => {
  const { internalItem, item } = slotProps;
  const title = internalItem?.title;

  if (title) {
    return String(title);
  }

  if (item?.title) {
    return String(item.title);
  }

  const raw = item?.raw ?? item;

  if (isNil(raw)) {
    return '';
  }

  if (!isObject(raw)) {
    return String(raw);
  }

  const rawTitle = raw.name ?? raw.text ?? raw.title ?? raw.value;

  return isNil(rawTitle) || isObject(rawTitle) ? '' : String(rawTitle);
};
