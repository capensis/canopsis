import {
  computed,
  onBeforeUnmount,
  ref,
  unref,
  watch,
} from 'vue';

/** Default number of items appended per infinite scroll page. */
export const INFINITE_SCROLL_DEFAULT_LIMIT = 20;

/** Root margin applied when observing the append item element. */
export const INFINITE_SCROLL_APPEND_ITEM_ROOT_MARGIN = '50px 0px';

/**
 * Client-side infinite scroll that exposes a growing window over a source list.
 *
 * `items` is a computed slice `source.slice(0, loadedCount)`, so it reacts to source mutations
 * (e.g. removing an item) without resetting pagination or scroll position: only the changed rows
 * are re-rendered while the loaded window size is preserved.
 * Observes an append item element at the bottom of the scroll container to load more automatically.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Array>|Array} options.sourceItems - Full source list to paginate.
 * @param {number} [options.limit=INFINITE_SCROLL_DEFAULT_LIMIT] - Number of items appended per page.
 * @param {import('vue').Ref<boolean>|boolean} [options.isLoading] - Blocks loading more while pending.
 * @param {string} [options.appendItemRootMargin=INFINITE_SCROLL_APPEND_ITEM_ROOT_MARGIN]
 * Root margin for append item observation.
 * @returns {{
 *   items: import('vue').ComputedRef<Array>,
 *   hasMore: import('vue').ComputedRef<boolean>,
 *   loadMore: function,
 *   scrollContainerElement: import('vue').Ref<HTMLElement|undefined>,
 *   appendItemElement: import('vue').Ref<HTMLElement|undefined>,
 * }}
 */
export const useInfiniteScroll = ({
  sourceItems,
  limit = INFINITE_SCROLL_DEFAULT_LIMIT,
  isLoading,
  appendItemRootMargin = INFINITE_SCROLL_APPEND_ITEM_ROOT_MARGIN,
} = {}) => {
  const loadedCount = ref(unref(limit));
  const scrollContainerElement = ref(null);
  const appendItemElement = ref(null);

  let appendObserver = null;

  const items = computed(() => (unref(sourceItems) ?? []).slice(0, loadedCount.value));
  const hasMore = computed(() => (unref(sourceItems) ?? []).length > items.value.length);

  /**
   * Grows the loaded window by one page.
   */
  const loadMore = () => {
    if (!hasMore.value) {
      return;
    }

    loadedCount.value += unref(limit);
  };

  /**
   * Disconnects the append item intersection observer.
   */
  const disconnectAppendObserver = () => {
    appendObserver?.disconnect();
    appendObserver = null;
  };

  /**
   * Handles append item intersection and loads the next page when needed.
   *
   * @param {IntersectionObserverEntry[]} entries
   */
  const handleAppendItemIntersection = (entries) => {
    const [entry] = entries;

    if (!entry?.isIntersecting || unref(isLoading) || !hasMore.value) {
      return;
    }

    loadMore();
  };

  /**
   * Creates the append item intersection observer and attaches it to the scroll container.
   */
  const connectAppendObserver = () => {
    disconnectAppendObserver();

    if (!scrollContainerElement.value) {
      return;
    }

    appendObserver = new IntersectionObserver(handleAppendItemIntersection, {
      root: scrollContainerElement.value,
      rootMargin: unref(appendItemRootMargin),
      threshold: 0,
    });

    if (appendItemElement.value) {
      appendObserver.observe(appendItemElement.value);
    }
  };

  watch(scrollContainerElement, connectAppendObserver);

  watch(appendItemElement, () => {
    if (appendItemElement.value && appendObserver) {
      appendObserver.observe(appendItemElement.value);
    }
  });

  onBeforeUnmount(disconnectAppendObserver);

  return {
    items,
    hasMore,
    loadMore,
    scrollContainerElement,
    appendItemElement,
  };
};
