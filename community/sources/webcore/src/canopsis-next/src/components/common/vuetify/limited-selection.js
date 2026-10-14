import { computed, useAttrs, useSlots } from 'vue';

import { getSlotItemTitle } from './slot-item';

/**
 * @typedef {Object} LimitedSelectionProps
 * Props consumed by `useLimitedSelection`.
 *
 * @property {number} [showCount] Max visible chips/selections before the rest are collapsed into a count.
 */

/**
 * @typedef {Object} LimitedSelectionOptions
 * Options for `useLimitedSelection`.
 *
 * @property {boolean} [hideSingleSelectionSlot=false] When `true`, do not forward `#selection` on single
 *   (non-multiple, non-chips) autocomplete/combobox fields so Vuetify keeps the value editable in the input.
 */

/**
 * @callback IsForwardedSlot
 * Whether a parent slot should be forwarded to the underlying Vuetify field.
 *
 * @param {string} name Slot name.
 * @returns {boolean}
 */

/**
 * @callback GetForwardedSlotNames
 * Parent slot names that should be forwarded on the current render.
 *
 * @returns {string[]}
 */

/**
 * @typedef {Object} LimitedSelectionReturn
 * Bindings returned by `useLimitedSelection`.
 *
 * @property {import('vue').ComputedRef<boolean>} hasLimitedSelection Whether extra selections should be collapsed
 *   with `showCount`.
 * @property {import('vue').ComputedRef<boolean>} isMultiple Whether the field is multiple or uses chips.
 * @property {import('vue').ComputedRef<number>} selectedCount Number of selected values when the model is an array.
 * @property {GetForwardedSlotNames} getForwardedSlotNames Slot names to forward (evaluated during render so new
 *   parent slots are picked up).
 * @property {typeof getSlotItemTitle} getSelectionTitle Resolves the display title from a selection/chip slot payload.
 */

/**
 * Props for collapsing extra selected items into a count.
 *
 * @type {LimitedSelectionProps}
 */
export const limitedSelectionProps = {
  showCount: {
    type: Number,
    default: undefined,
  },
};

/**
 * Handles limited chip/selection rendering and which parent slots are forwarded to Vuetify fields.
 *
 * `#chip` / `#selection` are not forwarded when `showCount` is set (the composable owns those slots).
 * With `hideSingleSelectionSlot`, `#selection` is also dropped for single fields so the selected value stays
 * in the input and can be edited, matching Vuetify 2 autocomplete behavior.
 *
 * @param {LimitedSelectionProps} props Field props, including optional `showCount`.
 * @param {import('vue').Ref<*>} inputValue Current field model.
 * @param {LimitedSelectionOptions} [options] Slot-forwarding options.
 * @returns {LimitedSelectionReturn}
 */
export const useLimitedSelection = (props, inputValue, { hideSingleSelectionSlot = false } = {}) => {
  const slots = useSlots();
  const attrs = useAttrs();

  const hasLimitedSelection = computed(() => props.showCount > 0);
  const isMultiple = computed(() => Boolean(attrs.multiple) || Boolean(attrs.chips));
  const selectedCount = computed(() => (
    Array.isArray(inputValue.value) ? inputValue.value.length : 0
  ));
  /**
   * @type {IsForwardedSlot}
   */
  const isForwardedSlot = (name) => {
    if (hasLimitedSelection.value && (name === 'chip' || name === 'selection')) {
      return false;
    }

    if (hideSingleSelectionSlot && !isMultiple.value && name === 'selection') {
      return false;
    }

    return true;
  };

  const getForwardedSlotNames = () => Object.keys(slots).filter(isForwardedSlot);

  const getSelectionTitle = getSlotItemTitle;

  return {
    hasLimitedSelection,
    isMultiple,
    selectedCount,
    getForwardedSlotNames,
    getSelectionTitle,
  };
};
