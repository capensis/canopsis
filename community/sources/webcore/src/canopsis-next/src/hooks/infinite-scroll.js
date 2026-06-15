import {
  computed,
  nextTick,
  onBeforeUnmount,
  ref,
  unref,
  watch,
} from 'vue';

/** Default number of items appended per infinite scroll page. */
export const INFINITE_SCROLL_DEFAULT_LIMIT = 20;

/**
 * Root margin applied when observing the append item element.
 *
 * Loads the next page well before the bottom sentinel reaches the viewport so fast scrolling does
 * not run past the last appended item into empty space.
 */
export const INFINITE_SCROLL_APPEND_ITEM_ROOT_MARGIN = '50px 0px';

/**
 * Builds a stable signature from source item ids to detect list changes without a deep watch.
 *
 * @param {Array<{ _id?: string|number, id?: string|number }>} [source]
 * @returns {string}
 */
const getSourceSignature = (source = []) => source.map(({ _id, id }) => _id ?? id).join(',');

/**
 * Client-side infinite scroll that appends items from a source list page by page.
 *
 * Keeps a local `items` ref and pushes the next slice from `sourceItems` on each `loadMore` call.
 * Observes an append item element at the bottom of the scroll container to load more automatically.
 * Resets and reloads the first page when the source list signature changes.
 *
 * @param {Object} [options]
 * @param {import('vue').Ref<Array>|Array} options.sourceItems - Full source list to paginate.
 * @param {number} [options.limit=INFINITE_SCROLL_DEFAULT_LIMIT] - Number of items appended per page.
 * @param {import('vue').Ref<boolean>|boolean} [options.isLoading] - Blocks loading more while pending.
 * @param {string} [options.appendItemRootMargin=INFINITE_SCROLL_APPEND_ITEM_ROOT_MARGIN]
 * Root margin for append item observation.
 * @param {function} [options.onReset] - Called before the list is cleared on source change.
 * @returns {{
 *   items: import('vue').Ref<Array>,
 *   hasMore: import('vue').ComputedRef<boolean>,
 *   loadMore: function,
 *   resetItems: function,
 *   scrollContainerElement: import('vue').Ref<HTMLElement|undefined>,
 *   appendItemElement: import('vue').Ref<HTMLElement|undefined>,
 * }}
 */
export const useInfiniteScroll = ({
  sourceItems,
  limit = INFINITE_SCROLL_DEFAULT_LIMIT,
  isLoading,
  appendItemRootMargin = INFINITE_SCROLL_APPEND_ITEM_ROOT_MARGIN,
  onReset,
} = {}) => {
  const items = ref([]);
  const page = ref(0);
  const scrollContainerElement = ref(null);
  const appendItemElement = ref(null);

  let appendObserver = null;

  const hasMore = computed(() => (unref(sourceItems) ?? []).length > items.value.length);

  /**
   * Appends the next page of items from the source list.
   */
  const loadMore = () => {
    if (!hasMore.value) {
      return;
    }

    const source = unref(sourceItems) ?? [];
    const nextPage = page.value + 1;
    const start = page.value * unref(limit);

    items.value.push(...source.slice(start, nextPage * unref(limit)));
    page.value = nextPage;
  };

  /**
   * Clears appended items, resets pagination and reloads the first page.
   */
  const resetItems = () => {
    onReset?.();
    items.value = [];
    page.value = 0;
    loadMore();
  };

  /**
   * Disconnects the append item intersection observer.
   */
  const disconnectAppendObserver = () => {
    appendObserver?.disconnect();
    appendObserver = null;
  };

  /**
   * Re-observes the append item to force a fresh intersection check.
   *
   * IntersectionObserver only emits on a state change, so when a short page (e.g. the last one)
   * is appended without pushing the sentinel out of the root margin, no new entry is delivered.
   * Re-observing after the DOM updates re-triggers loading while the sentinel stays in the zone.
   */
  const reobserveAppendItem = () => {
    if (!appendObserver || !appendItemElement.value || !hasMore.value) {
      return;
    }

    appendObserver.unobserve(appendItemElement.value);
    appendObserver.observe(appendItemElement.value);
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
    nextTick(reobserveAppendItem);
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

  watch(() => getSourceSignature(unref(sourceItems)), resetItems, { immediate: true });

  onBeforeUnmount(disconnectAppendObserver);

  return {
    items,
    hasMore,
    loadMore,
    resetItems,
    scrollContainerElement,
    appendItemElement,
  };
};
