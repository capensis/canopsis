import {
  computed,
  ref,
  unref,
  watch,
  onMounted,
} from 'vue';
import { keyBy, pick, isArray, isString } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { mapIds } from '@/helpers/array';

import { usePendingHandler } from '@/hooks/query/pending';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useModelField } from '@/hooks/form';

/**
 * Custom hook for lazy searching and managing selected items with pagination and search capabilities.
 *
 * This hook integrates with the Vue composition API and organization-specific utilities to provide a reactive
 * interface for fetching, selecting, and managing items. It supports lazy loading, search, and pagination.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Ref|Array} options.value - The current value or list of values to be managed.
 * @param {string} options.isKey - The key used to identify each item uniquely.
 * @param {string} options.idParamsKey - The key used for query parameters when fetching selected items.
 * @param {number} [options.limit = PAGINATION_LIMIT] - The limit for pagination, determining how many items to
 * fetch per page.
 * @param {Function} options.fetchHandler - The asynchronous function used to fetch data from the server.
 * @param {boolean} options.addable - The flag for indicating possibility to add new item.
 * @param {Function} emit - The emit function for Vue events, used to update the model.
 * @returns {Object} An object containing methods and properties for managing search and selection:
 * - `selectedItems`: {Ref<Array>} A reactive reference to the currently selected items.
 * - `items`: {ComputedRef<Array>} A computed reference to the list of items fetched and managed by the hook.
 * - `wholePending`: {ComputedRef<boolean>} A computed reference indicating if any fetch operation is pending.
 * - `hasMoreItems`: {ComputedRef<boolean>} A computed reference indicating if there are more items to fetch.
 * - `fetchItems`: {Function} A function to fetch items based on the current query state.
 * - `fetchMoreItems`: {Function} A function to fetch the next page of items.
 * - `changeSelectedItems`: {Function} A function to update the selected items and emit changes.
 * - `updateSearch`: {Function} A function to update the search query and trigger a fetch.
 */
export const useLazySearch = ({ value, idKey, idParamsKey, limit = PAGINATION_LIMIT, fetchHandler, addable }, emit) => {
  const pageCount = ref(1);
  const itemsByValue = ref({});
  const selectedItems = ref([]);

  const { updateModel } = useModelField({}, emit);

  /**
   * Computed property to get the list of items from the itemsByValue map.
   * @type {ComputedRef<Array>}
   */
  const items = computed(() => Object.values(itemsByValue.value));

  /**
   * Computed property to convert the value into an array format.
   * @type {ComputedRef<Array>}
   */
  const arrayValue = computed(() => {
    const unwrappedValue = unref(value);

    if (isArray(unwrappedValue)) {
      return unwrappedValue;
    }

    return unwrappedValue ? [unwrappedValue] : [];
  });

  /**
   * FETCH VALUE ITEMS
   */
  const {
    pending: valuesPending,
    handler: initializeSelectedItems,
  } = usePendingHandler(async () => {
    const selectedItemsFromItemsByValue = arrayValue.value.map(item => itemsByValue.value[item]).filter(Boolean);

    if (selectedItemsFromItemsByValue.length === arrayValue.value.length) {
      selectedItems.value = selectedItemsFromItemsByValue;

      return;
    }

    if (!arrayValue.value.length) {
      return;
    }

    const { data } = await fetchHandler({
      params: {
        limit: arrayValue.value.length,
        [unref(idParamsKey)]: arrayValue.value,
      },
    });

    const unwrappedIdKey = unref(idKey);
    const dataById = keyBy(data, unwrappedIdKey);

    selectedItems.value = arrayValue.value.map(item => (
      dataById[item[unwrappedIdKey] || item] ?? ({ [unwrappedIdKey]: item })
    ));
  }, true);

  /**
   * MAIN FETCH LOGIC
   */
  const {
    pending,
    query,
    fetchHandlerWithQuery: fetchItems,
    updateQueryPage,
    updateQuerySearch,
  } = usePendingWithLocalQuery({
    initialQuery: {
      page: 1,
      limit: unref(limit),
      search: '',
    },
    fetchHandler: async (params) => {
      const { data, meta } = await fetchHandler({
        params,
      });

      pageCount.value = meta.page_count;

      itemsByValue.value = {
        ...(params.page !== 1 ? itemsByValue.value : {}),
        ...keyBy(data, unref(idKey)),
        ...pick(itemsByValue.value, arrayValue.value),
      };
    },
  });

  /**
   * Computed property to determine if there are more items to fetch.
   * @type {ComputedRef<boolean>}
   */
  const hasMoreItems = computed(() => pageCount.value > query.value.page);

  /**
   * Computed property to determine if any fetch operation is pending.
   * @type {ComputedRef<boolean>}
   */
  const wholePending = computed(() => pending.value || valuesPending.value);

  /**
   * Function to fetch the next page of items.
   */
  const fetchMoreItems = () => updateQueryPage(query.value.page + 1);

  /**
   * Function to update the selected items and emit changes.
   * @param {Array} newSelectedTags - The new list of selected tags.
   */
  const changeSelectedItems = (newSelectedTags) => {
    const unwrappedIdKey = unref(idKey);
    const unwrappedAddable = unref(addable);

    selectedItems.value = (
      unwrappedAddable
        ? newSelectedTags
        : newSelectedTags.filter(tag => !isString(tag))
    ).map(tag => (tag[unwrappedIdKey] ? tag : { [unwrappedIdKey]: tag }));

    updateModel(mapIds(selectedItems.value, unwrappedIdKey));
  };

  /**
   * Removes an item from the `selectedItems` array by its index and updates the selection.
   *
   * This function filters out the item at the specified index from the `selectedItems` array
   * and then updates the selection by calling `changeSelectedItems`. It is useful for managing
   * the removal of items in a list where the index is known.
   *
   * @param {number} index - The index of the item to be removed from the `selectedItems` array.
   */
  const removeItemFromSelectedItemsByIndex = index => (
    changeSelectedItems(selectedItems.value.filter((item, itemIndex) => itemIndex !== index))
  );

  watch(value, () => initializeSelectedItems());

  onMounted(() => {
    if (idParamsKey) {
      initializeSelectedItems();
    }

    fetchItems();
  });

  return {
    selectedItems,
    items,
    wholePending,
    hasMoreItems,
    fetchItems,
    fetchMoreItems,
    changeSelectedItems,
    removeItemFromSelectedItemsByIndex,
    updateSearch: updateQuerySearch,
  };
};
