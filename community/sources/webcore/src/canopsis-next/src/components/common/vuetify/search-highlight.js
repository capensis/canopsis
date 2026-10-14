import { computed, ref, watch } from 'vue';

import { getSlotItemTitle } from './slot-item';

/**
 * @typedef {Object<string, *>} SearchHighlightAttrs
 * Fallthrough attrs from an autocomplete/combobox wrapper.
 *
 * @property {string} [search] Controlled search string.
 * @property {Function} [onUpdate:search] Parent listener for `update:search`.
 */

/**
 * @typedef {Object} SearchHighlightItemProps
 * Vuetify list-item bindings from an `#item` slot payload, without `title` (the title is rendered separately
 * so it can be highlighted).
 */

/**
 * @callback GetHighlightedItemProps
 * List-item props for a highlighted `#item` slot, with `title` stripped so the wrapper can render it.
 *
 * @param {Object} [slotProps] Vuetify `#item` slot payload.
 * @param {Object} [slotProps.props] Bindings meant for `v-list-item`.
 * @returns {SearchHighlightItemProps|{}}
 */

/**
 * @typedef {Object} SearchHighlightReturn
 * Bindings returned by `useSearchHighlight`.
 *
 * @property {import('vue').Ref<string>} searchQuery Current search string used to mask/highlight item titles.
 * @property {Function} onSearchUpdate Syncs `searchQuery` and re-emits `update:search` to the parent.
 * @property {import('vue').ComputedRef<SearchHighlightAttrs>} fieldAttrsWithoutSearchListener Attrs to bind on the
 *   Vuetify field, with `onUpdate:search` removed so the composable owns that listener.
 * @property {typeof getSlotItemTitle} getHighlightedText Resolves the item title to highlight from a slot payload.
 * @property {GetHighlightedItemProps} getHighlightedItemProps `v-list-item` props without `title`.
 */

/**
 * Tracks autocomplete/combobox search and prepares item titles for highlight masking.
 *
 * The wrapper listens to `update:search` itself so `searchQuery` stays in sync for `v-list-item-mask`.
 * `onUpdate:search` is stripped from fallthrough attrs to avoid double-handling.
 *
 * @param {SearchHighlightAttrs} attrs Component fallthrough attrs (`useAttrs()`).
 * @returns {SearchHighlightReturn}
 */
export const useSearchHighlight = (attrs) => {
  const searchQuery = ref(typeof attrs.search === 'string' ? attrs.search : '');

  watch(
    () => attrs.search,
    (value) => {
      if (value !== undefined) {
        searchQuery.value = value ?? '';
      }
    },
  );

  const onSearchUpdate = (value) => {
    searchQuery.value = value ?? '';
    attrs['onUpdate:search']?.(value);
  };

  const fieldAttrsWithoutSearchListener = computed(() => {
    const nextAttrs = { ...attrs };

    delete nextAttrs['onUpdate:search'];

    return nextAttrs;
  });

  const getHighlightedText = getSlotItemTitle;

  /**
   * @type {GetHighlightedItemProps}
   */
  const getHighlightedItemProps = (slotProps) => {
    const itemProps = slotProps?.props;

    if (!itemProps) {
      return {};
    }

    const { title, ...rest } = itemProps;

    return rest;
  };

  return {
    searchQuery,
    onSearchUpdate,
    fieldAttrsWithoutSearchListener,
    getHighlightedText,
    getHighlightedItemProps,
  };
};
