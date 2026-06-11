import { computed, unref } from 'vue';

import { advancedSearchToForm } from '@/helpers/search/advanced-search';

import { useInfiniteScroll } from '@/hooks/infinite-scroll';
import { useIntersectionVirtualScroll } from '@/hooks/intersection-virtual-scroll';

export const ADVANCED_SEARCH_HISTORY_ITEM_MIN_HEIGHT = 48;

/**
 * Prepares advanced search history items with append-based infinite scroll and intersection virtual scroll.
 *
 * @param {import('vue').Ref<Array>|Array} searches - Raw saved search items from the parent.
 * @returns {{
 *   items: import('vue').Ref<Array>,
 *   scrollContainerElement: import('vue').Ref<HTMLElement|undefined>,
 *   appendItemElement: import('vue').Ref<HTMLElement|undefined>,
 *   setItemRef: function(string|number): function(HTMLElement|undefined),
 *   isItemVisible: function(string|number): boolean,
 *   getItemPlaceholderStyle: function(string|number): Object,
 *   listMinWidth: import('vue').Ref<number>,
 *   updateListWidth: function,
 * }}
 */
export const useAdvancedSearchHistory = (searches) => {
  const preparedSearches = computed(() => (unref(searches) ?? []).map(search => ({
    _id: search._id,
    pinned: search.pinned,
    rules: advancedSearchToForm(search),
  })));

  let resetVirtualScroll = () => {};

  const {
    items,
    scrollContainerElement,
    appendItemElement,
  } = useInfiniteScroll({
    sourceItems: preparedSearches,
    onReset: () => resetVirtualScroll(),
  });

  const {
    resetVirtualScroll: resetVirtualScrollHandler,
    setItemRef,
    isItemVisible,
    getItemPlaceholderStyle,
    listMinWidth,
    updateListWidth,
  } = useIntersectionVirtualScroll({
    scrollContainerElement,
    itemMinHeight: ADVANCED_SEARCH_HISTORY_ITEM_MIN_HEIGHT,
  });

  resetVirtualScroll = resetVirtualScrollHandler;

  return {
    items,
    scrollContainerElement,
    appendItemElement,
    setItemRef,
    isItemVisible,
    getItemPlaceholderStyle,
    listMinWidth,
    updateListWidth,
  };
};
