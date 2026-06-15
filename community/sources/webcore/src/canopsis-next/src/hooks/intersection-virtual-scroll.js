import { onBeforeUnmount, ref, unref, watch } from 'vue';

/**
 * Root margin applied when observing list item visibility.
 *
 * A large vertical buffer pre-renders rows well before they enter the viewport so fast scrolling
 * does not outrun the asynchronous IntersectionObserver callbacks (which would flash placeholders).
 */
export const INTERSECTION_VIRTUAL_SCROLL_ITEM_ROOT_MARGIN = '600px 0px';

/** Default placeholder height in pixels when an item height has not been measured yet. */
export const INTERSECTION_VIRTUAL_SCROLL_DEFAULT_ITEM_MIN_HEIGHT = 48;

/**
 * Manual virtual scroll based on Intersection Observer.
 *
 * Renders only items that intersect the scroll container and replaces off-screen rows with
 * placeholders that preserve measured height and width.
 *
 * @param {Object} options
 * @param {import('vue').Ref<HTMLElement|undefined>} options.scrollContainerElement
 * Scroll container used as observer root.
 * @param {string} [options.itemRootMargin=INTERSECTION_VIRTUAL_SCROLL_ITEM_ROOT_MARGIN]
 * Root margin for item visibility observation.
 * @param {number} [options.itemMinHeight=INTERSECTION_VIRTUAL_SCROLL_DEFAULT_ITEM_MIN_HEIGHT]
 * Fallback placeholder height in pixels.
 * @returns {{
 *   setItemRef: function(string|number): function(HTMLElement|undefined),
 *   isItemVisible: function(string|number): boolean,
 *   getItemPlaceholderStyle: function(string|number): Object,
 *   listMinWidth: import('vue').Ref<number>,
 *   updateListWidth: function(HTMLElement=),
 *   resetVirtualScroll: function,
 * }}
 */
export const useIntersectionVirtualScroll = ({
  scrollContainerElement,
  itemRootMargin = INTERSECTION_VIRTUAL_SCROLL_ITEM_ROOT_MARGIN,
  itemMinHeight = INTERSECTION_VIRTUAL_SCROLL_DEFAULT_ITEM_MIN_HEIGHT,
} = {}) => {
  const visibilityById = ref({});
  const itemHeights = ref({});
  const itemWidths = ref({});
  const listMinWidth = ref(350);

  const observedElements = new Map();
  const itemRefSetters = new Map();
  let itemObserver = null;

  /**
   * Disconnects the item intersection observer.
   */
  const disconnectObserver = () => {
    itemObserver?.disconnect();
    itemObserver = null;
  };

  /**
   * Clears visibility, size caches and unobserves all registered item elements.
   */
  const resetVisibilityState = () => {
    visibilityById.value = {};
    itemHeights.value = {};
    itemWidths.value = {};
    listMinWidth.value = 0;
    observedElements.forEach((element) => {
      itemObserver?.unobserve(element);
    });
    observedElements.clear();
    itemRefSetters.clear();
  };

  /**
   * Updates the minimum list width from the widest measured container or item layout.
   *
   * @param {HTMLElement} [element]
   */
  const updateListWidth = (element = unref(scrollContainerElement)) => {
    if (!element) {
      return;
    }

    const containerWidth = Math.max(element.scrollWidth, element.offsetWidth);

    if (containerWidth > listMinWidth.value) {
      listMinWidth.value = containerWidth;
    }
  };

  /**
   * Handles item visibility intersection changes.
   *
   * Batches every entry of a single observer callback into at most one clone per
   * reactive object, instead of cloning on each entry.
   *
   * @param {IntersectionObserverEntry[]} entries
   */
  const handleItemIntersection = (entries) => {
    let nextVisibility;
    let nextHeights;
    let nextWidths;
    let maxItemWidth = listMinWidth.value;
    let hasVisibleItem = false;

    entries.forEach((entry) => {
      const element = entry.target;
      const { virtualScrollId: id } = element.dataset;

      if (!id) {
        return;
      }

      const isVisible = entry.isIntersecting;
      const { offsetWidth } = element;

      if (offsetWidth) {
        nextWidths = nextWidths ?? { ...itemWidths.value };
        nextWidths[id] = offsetWidth;

        const layoutWidth = Math.max(element.scrollWidth, offsetWidth);

        if (layoutWidth > maxItemWidth) {
          maxItemWidth = layoutWidth;
        }
      }

      if (isVisible) {
        hasVisibleItem = true;
      } else if (element.offsetHeight) {
        nextHeights = nextHeights ?? { ...itemHeights.value };
        nextHeights[id] = element.offsetHeight;
      }

      if (!!visibilityById.value[id] !== isVisible) {
        nextVisibility = nextVisibility ?? { ...visibilityById.value };
        nextVisibility[id] = isVisible;
      }
    });

    if (nextWidths) {
      itemWidths.value = nextWidths;
    }

    if (nextHeights) {
      itemHeights.value = nextHeights;
    }

    if (nextVisibility) {
      visibilityById.value = nextVisibility;
    }

    if (maxItemWidth > listMinWidth.value) {
      listMinWidth.value = maxItemWidth;
    }

    if (hasVisibleItem) {
      updateListWidth();
    }
  };

  /**
   * Creates the item intersection observer and attaches it to the scroll container.
   */
  const connectObserver = () => {
    disconnectObserver();

    if (!unref(scrollContainerElement)) {
      return;
    }

    itemObserver = new IntersectionObserver(handleItemIntersection, {
      root: unref(scrollContainerElement),
      rootMargin: unref(itemRootMargin),
      threshold: 0,
    });

    observedElements.forEach((element) => {
      itemObserver.observe(element);
    });
  };

  /**
   * Returns a stable ref callback that registers an item element with the item intersection
   * observer. The callback is memoized per id so Vue keeps the same ref identity across re-renders
   * and does not unobserve/re-observe every element on each visibility change.
   *
   * @param {string|number} id
   * @returns {function(HTMLElement|undefined)}
   */
  const setItemRef = (id) => {
    if (!itemRefSetters.has(id)) {
      itemRefSetters.set(id, (element) => {
        const previousElement = observedElements.get(id);

        if (previousElement && previousElement !== element) {
          itemObserver?.unobserve(previousElement);
          observedElements.delete(id);
        }

        if (!element) {
          return;
        }

        // eslint-disable-next-line no-param-reassign
        element.dataset.virtualScrollId = String(id);
        observedElements.set(id, element);
        itemObserver?.observe(element);
      });
    }

    return itemRefSetters.get(id);
  };

  /**
   * Checks whether the item with the given id is currently visible in the scroll container.
   *
   * @param {string|number} id
   * @returns {boolean}
   */
  const isItemVisible = id => !!visibilityById.value[id];

  /**
   * Returns inline styles for an off-screen item placeholder.
   *
   * @param {string|number} id
   * @returns {Object}
   */
  const getItemPlaceholderStyle = (id) => {
    const height = itemHeights.value[id] ?? unref(itemMinHeight);
    const width = itemWidths.value[id] ?? listMinWidth.value;

    return {
      minHeight: `${height}px`,
      ...(width ? { minWidth: `${width}px` } : {}),
    };
  };

  watch(() => unref(scrollContainerElement), () => {
    connectObserver();
    updateListWidth();
  });

  onBeforeUnmount(disconnectObserver);

  return {
    setItemRef,
    isItemVisible,
    getItemPlaceholderStyle,
    listMinWidth,
    updateListWidth,
    resetVirtualScroll: resetVisibilityState,
  };
};
