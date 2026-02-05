import {
  keyBy,
  pick,
  isArray,
  isString,
  isNumber,
  isUndefined,
  isEmpty,
  isObject,
  debounce,
  uniq,
  uniqBy,
} from 'lodash';
import {
  computed,
  ref,
  unref,
  watch,
  set,
  onMounted,
} from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { deepKeyBy, mapIds } from '@/helpers/array';

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
 * @param {string} options.childrenKey - The key used for children checking.
 * @param {string} options.idParamsKey - The key used for query parameters when fetching selected items.
 * @param {Function} options.fetchHandler - The asynchronous function used to fetch data from the server.
 * @param {boolean} options.addable - The flag for indicating possibility to add new item.
 * @param {boolean} options.multiple - The flag for indicating possibility to choose more than one item.
 * @param {boolean} options.returnObject - The flag for indicating possibility to return object instead of id.
 * @param {boolean} options.attachValue - The flag for indicating possibility to attach the value to the items.
 * @param {Object} options.initialQuery - The initial query object.
 * @param {number} options.delay - The delay number value for debouncing.
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
export const useLazySearch = ({
  value,
  idKey,
  idParamsKey,
  childrenKey = 'items',
  fetchHandler,
  addable,
  multiple,
  returnObject,
  attachValue,
  initialQuery = { page: 1, limit: PAGINATION_LIMIT, search: '' },
  delay = 100,
}, emit) => {
  const pageCount = ref(1);
  const itemsByValue = ref({});
  const itemsValues = ref([]);
  const selectedItems = ref([]);

  const { updateModel } = useModelField({}, emit);

  /**
   * Computed property to get the list of items from the itemsByValue map.
   * @type {ComputedRef<Array>}
   */
  const items = computed(() => itemsValues.value.map(itemValue => itemsByValue.value[itemValue]));

  /**
   * Computed property to convert the value into an array format.
   * @type {ComputedRef<Array>}
   */
  const arrayValue = computed(() => {
    const unwrappedValue = unref(value);

    if ((!unwrappedValue || isEmpty(unwrappedValue)) && !isNumber(unwrappedValue)) {
      return [];
    }

    if (isArray(unwrappedValue)) {
      return unwrappedValue;
    }

    return [unwrappedValue];
  });

  /**
   * FETCH VALUE ITEMS
   */
  const {
    pending: valuesPending,
    handler: initializeSelectedItems,
  } = usePendingHandler(async () => {
    if (!arrayValue.value.length && selectedItems.value.length) {
      selectedItems.value = [];
      return;
    }

    const selectedItemsFromItemsByValue = arrayValue.value.map(item => itemsByValue.value[item]).filter(Boolean);

    if (selectedItemsFromItemsByValue.length === arrayValue.value.length) {
      selectedItems.value = selectedItemsFromItemsByValue;

      return;
    }

    const { data } = await unref(fetchHandler)({
      params: {
        limit: arrayValue.value.length,
        [unref(idParamsKey)]: arrayValue.value,
      },
    });
    const unwrappedIdKey = unref(idKey);
    const unwrappedChildrenKey = unref(childrenKey);

    const dataById = deepKeyBy(data, unwrappedIdKey, unwrappedChildrenKey);

    selectedItems.value = arrayValue.value.map(item => (
      dataById[item[unwrappedIdKey] ?? item]
        ?? (isObject(item) ? item : ({ [unwrappedIdKey]: item, noData: true }))
    ));

    if (attachValue) {
      selectedItems.value.forEach((item) => {
        const id = item[unwrappedIdKey];

        if (id && !itemsByValue.value[id]) {
          set(itemsByValue.value, id, item);
        }
      });
    }
  }, true);

  /**
     * MAIN FETCH LOGIC
     */
  const {
    pending,
    query,
    fetchHandlerWithQuery: fetchItems,
    updateQuery,
    updateQueryPage,
  } = usePendingWithLocalQuery({
    initialQuery: unref(initialQuery),
    fetchHandler: async (params) => {
      const unwrappedIdKey = unref(idKey);

      const { data, meta } = await unref(fetchHandler)({
        params,
      });

      pageCount.value = meta.page_count;

      /**
      * We need to use it for saving order of items
      */
      itemsValues.value = uniq([
        ...(params.page && params.page !== 1 ? itemsValues.value : []),
        ...data.map(item => item[unwrappedIdKey]),
      ]);

      itemsByValue.value = {
        ...(params.page && params.page !== 1 ? itemsByValue.value : {}),
        ...keyBy(data, unwrappedIdKey),
        ...pick(itemsByValue.value, arrayValue.value),
      };
    },
  });

  /**
   * Update search field in query with page updating
   *
   * @param {string} [newSearch = '']
   */
  const updateSearch = debounce((newSearch = '') => {
    if (query.value.search !== newSearch) {
      updateQuery({
        search: newSearch,
        page: 1,
      });
    }
  }, unref(delay));

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
   * Computed property for first selected item.
   * @type {ComputedRef<Object | undefined>}
   */
  const selectedItem = computed(() => selectedItems.value?.[0]);

  /**
   * Function to fetch the next page of items.
   */
  const fetchMoreItems = () => updateQueryPage(query.value.page + 1);

  /**
   * Function to update the selected items and emit changes.
   * @param {Array} newSelectedItems - The new list of selected items.
   */
  const changeSelectedItems = (newSelectedItems) => {
    if (!newSelectedItems) {
      selectedItems.value = [];

      updateModel('');

      return;
    }

    const unwrappedIdKey = unref(idKey);
    const unwrappedAddable = unref(addable);
    const unwrappedMultiple = unref(multiple);

    let preparedNewSelectedItems;

    if (!isArray(newSelectedItems)) {
      preparedNewSelectedItems = [newSelectedItems];
    } else {
      preparedNewSelectedItems = unwrappedMultiple ? newSelectedItems : newSelectedItems.slice(-1);
    }

    selectedItems.value = uniqBy(
      (
        unwrappedAddable
          ? preparedNewSelectedItems
          : preparedNewSelectedItems.filter(item => !isString(item))
      ).map(item => (
        isUndefined(item[unwrappedIdKey])
          ? { [unwrappedIdKey]: item, noData: true }
          : item)),
      unwrappedIdKey,
    );

    if (returnObject) {
      updateModel(
        unwrappedMultiple
          ? selectedItems.value
          : selectedItems.value[0],
      );

      return;
    }

    updateModel(
      unwrappedMultiple
        ? mapIds(selectedItems.value, unwrappedIdKey)
        : selectedItems.value[0]?.[unwrappedIdKey],
    );
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
  watch(selectedItems, newSelectedItems => emit('update:selected-items', newSelectedItems));

  onMounted(() => {
    if (unref(idParamsKey)) {
      initializeSelectedItems();
    }

    fetchItems();
  });

  return {
    query,
    selectedItem,
    selectedItems,
    items,
    valuesPending,
    wholePending,
    hasMoreItems,
    fetchItems,
    fetchMoreItems,
    changeSelectedItems,
    removeItemFromSelectedItemsByIndex,
    updateSearch,
  };
};
