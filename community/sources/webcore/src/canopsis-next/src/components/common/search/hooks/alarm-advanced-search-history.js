import { computed, unref } from 'vue';

import { advancedSearchToForm } from '@/helpers/search/alarm-advanced-search';

import { useInfiniteScroll } from '@/hooks/infinite-scroll';
import { useIntersectionVirtualScroll } from '@/hooks/intersection-virtual-scroll';

export const ALARM_ADVANCED_SEARCH_HISTORY_ITEM_MIN_HEIGHT = 48;

/**
 * Prepares alarm advanced search history items with append-based infinite scroll and intersection virtual scroll.
 *
 * @param {import('vue').Ref<Array>|Array} searches - Raw saved search items from the parent.
 * @returns {{
 *   items: import('vue').ComputedRef<Array>,
 *   scrollContainerElement: import('vue').Ref<HTMLElement|undefined>,
 *   appendItemElement: import('vue').Ref<HTMLElement|undefined>,
 *   setItemRef: function(string|number): function(HTMLElement|undefined),
 *   isItemVisible: function(string|number): boolean,
 *   getItemPlaceholderStyle: function(string|number): Object,
 *   listMinWidth: import('vue').Ref<number>,
 *   updateListWidth: function,
 * }}
 */
export const useAlarmAdvancedSearchHistory = (searches) => {
  const preparedSearches = computed(() => (unref(searches) ?? []).map(search => ({
    _id: search._id,
    pinned: search.pinned,
    rules: advancedSearchToForm(search),
  })));

  const {
    items,
    scrollContainerElement,
    appendItemElement,
  } = useInfiniteScroll({
    sourceItems: preparedSearches,
  });

  const {
    setItemRef,
    isItemVisible,
    getItemPlaceholderStyle,
    listMinWidth,
    updateListWidth,
  } = useIntersectionVirtualScroll({
    scrollContainerElement,
    itemMinHeight: ALARM_ADVANCED_SEARCH_HISTORY_ITEM_MIN_HEIGHT,
  });

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
