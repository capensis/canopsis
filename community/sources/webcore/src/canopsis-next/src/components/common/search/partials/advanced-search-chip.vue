<template>
  <v-menu
    v-bind="menuProps"
    :value="opened"
    bottom
    @input="updateMenuOpened"
  >
    <template #activator="{ on }">
      <c-simple-tooltip :content="tooltip" :disabled="!tooltip" top>
        <template #activator="{ on: tooltipOn }">
          <span v-on="tooltipOn">
            <span v-if="multiple">
              <v-chip
                :close="closable && !disabled"
                :color="color"
                class="c-advanced-search__chip c-advanced-search__array-chip"
                v-on="multipleChipListeners"
              >
                <v-progress-circular
                  v-if="valuesPending"
                  color="primary"
                  size="16"
                  width="2"
                  indeterminate
                />
                <v-chip
                  v-for="(item, index) in selectedItems"
                  :key="item[itemValue]"
                  :close="!disabled"
                  @click:close="closeChildChip(index)"
                >
                  {{ item[itemText] ?? item[itemValue] }}
                </v-chip>
                <span v-if="!disabled">
                  <input
                    v-model="inputValue"
                    ref="inputElement"
                    v-bind="inputAttributes"
                    :type="inputType"
                    class="ml-1"
                    autocomplete="off"
                    @keydown="keydownInput"
                    @focus="focusInput"
                    @input="updateInputValue"
                  >
                </span>
              </v-chip>
            </span>
            <span v-else>
              <input
                v-show="active || alwaysActive"
                ref="inputElement"
                v-bind="inputAttributes"
                :value="inputValue"
                :placeholder="inputPlaceholder"
                :type="inputType"
                class="ml-1"
                autocomplete="off"
                v-on="on"
                @keydown="keydownInput"
                @mouseup="focusInput"
                @input="updateInputValue"
              >
              <v-chip
                v-if="!active && !alwaysActive"
                :color="color"
                :close="closable && !disabled"
                class="c-advanced-search__chip"
                v-on="chipListeners"
              >
                <v-progress-circular
                  v-if="valuesPending"
                  color="primary"
                  size="16"
                  width="2"
                  indeterminate
                />
                <span v-else>
                  <c-simple-tooltip v-if="icon" :content="icon.tooltip" top>
                    <template #activator="{ on: secondTooltipOn }">
                      <v-icon class="mr-2" small v-on="secondTooltipOn">{{ icon.icon }}</v-icon>
                    </template>
                  </c-simple-tooltip>
                  {{ chipText }}
                </span>
              </v-chip>
            </span>
          </span>
        </template>
      </c-simple-tooltip>
    </template>
    <variables-list
      :value="selectedItems"
      :items="lazyItems"
      :search="inputValue"
      :pending="fetchItems && pending"
      :item-text="itemText"
      :item-value="itemValue"
      :children-key="childrenKey"
      return-object
      @input="selectItem"
      @fetch:more="showMore"
    >
      <template v-if="hasAnyDisabledItem" #prepend="">
        <v-list-item class="font-italic grey--text">
          <v-list-item-content>
            <v-list-item-title>{{ $t('advancedSearch.listDisabledMessage') }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-divider />
      </template>
      <template #no-data="">
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title v-html="noDataText" />
          </v-list-item-content>
        </v-list-item>
      </template>
    </variables-list>
  </v-menu>
</template>

<script>
import { isArray } from 'lodash';
import {
  computed,
  ref,
  watch,
  toRef,
  inject,
  nextTick,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { KEY_CODES, REGISTER_LAST_INPUT_FOCUS_KEY } from '@/constants';

import { filterAdvancedSearchItems } from '@/helpers/search/advanced-search';

import { useI18n } from '@/hooks/i18n';
import { useLazySearch } from '@/hooks/form/lazy-search';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

export default {
  components: { VariablesList },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Number, Array, Boolean],
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    fetchItems: {
      type: Function,
      required: false,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    childrenKey: {
      type: String,
      default: 'items',
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    active: {
      type: Boolean,
      default: false,
    },
    allowText: {
      type: Boolean,
      default: false,
    },
    alwaysActive: {
      type: Boolean,
      default: false,
    },
    focusOnMount: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
    first: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    number: {
      type: Boolean,
      default: false,
    },
    tooltip: {
      type: String,
      default: '',
    },
    inputAttributes: {
      type: Object,
      required: false,
    },
    inputPreparer: {
      type: Function,
      required: false,
    },
  },
  setup(props, { emit }) {
    const focusRegister = inject(REGISTER_LAST_INPUT_FOCUS_KEY, {});

    const { t } = useI18n();

    const inputElement = ref(null);
    const inputValue = ref('');
    const opened = ref(false);

    /**
     * Fetches and filters items based on search parameters, returning the filtered data along with metadata about the
     * total count and page count.
     * We need it to simulate the same behavior as on real search request.
     *
     * @param {Object} [options = {}] - Options for fetching items.
     * @param {Object} [options.params = {}] - Parameters for filtering items.
     * @param {string} [options.params.search] - The search query used to filter items.
     * @returns {{ data: [], meta: { total_count: Number, page_count: 1 } }}
     */
    const defaultFetchItems = ({ params = {} } = {}) => {
      let data = props.items;

      if (params?.search) {
        const lowerSearch = String(params.search).toLowerCase();

        data = filterAdvancedSearchItems(data, item => item[props.itemText]?.toLowerCase().includes(lowerSearch));
      }

      return {
        data,
        meta: { total_count: data.length, page_count: 1 },
      };
    };

    const fetchHandler = computed(() => props.fetchItems ?? defaultFetchItems);

    const {
      selectedItems,
      updateSearch,
      valuesPending,
      hasMoreItems,
      items: lazyItems,
      wholePending: pending,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
    } = useLazySearch({
      fetchHandler,

      idParamsKey: 'ids',
      value: toRef(props, 'value'),
      addable: toRef(props, 'allowText'),
      idKey: toRef(props, 'itemValue'),
      multiple: toRef(props, 'multiple'),
    }, emit);

    const itemWithAlias = computed(() => selectedItems.value.find(({ alias } = {}) => alias));

    const menuProps = computed(() => ({
      openOnClick: false,
      disableKeys: true,
      closeOnContentClick: false,
      ignoreClickOutsideOnActivator: true,
      maxHeight: 304,
      nudgeBottom: 1,
      bottom: true,
      offsetY: true,
      transition: false,
      disabled: props.disabled,
    }));

    const inputType = computed(() => (props.number ? 'number' : 'text'));

    const hasAnyDisabledItem = computed(() => props.items.some(({ disabled }) => disabled));

    const inputPlaceholder = computed(() => (props.first ? t('advancedSearch.inputPlaceholder') : ''));

    const noDataText = computed(() => {
      let messageKey = 'common.noData';

      if (props.allowText) {
        messageKey = props.first ? 'advancedSearch.searchForThisText' : 'common.pressEnterToApply';
      }

      return t(messageKey);
    });

    const chipText = computed(() => selectedItems.value
      .map(item => item.chipText ?? item[props.itemText] ?? item[props.itemValue])
      .join(','));

    const hasNoDataFlagInSelectedItems = computed(() => selectedItems.value.some(({ noData }) => noData));

    const icon = computed(() => {
      if (itemWithAlias.value) {
        return {
          icon: 'alternate_email',
          tooltip: `infos.${itemWithAlias.value?.original?.name}`,
        };
      }

      if (props.first && hasNoDataFlagInSelectedItems.value) {
        return {
          icon: 'text_fields',
          tooltip: t('advancedSearch.searchByText'),
        };
      }

      return null;
    });

    /**
     * Updates the input value and triggers a search update with the new value.
     *
     * @param {string} newInputValue - The new value to be set for the input.
     */
    const setInputValue = (newInputValue, withSearch = true) => {
      inputValue.value = newInputValue;

      if (withSearch) {
        updateSearch(newInputValue);
      }
    };

    /**
     * Event handler that updates the input value based on the user's input event.
     * It extracts the value from the event's target and sets it as the new input value.
     *
     * @param {Event} event - The input event triggered by the user.
     */
    const updateInputValue = event => (
      setInputValue(props.inputPreparer ? props.inputPreparer(event.target.value) : event.target.value)
    );

    /**
     * Updates the state of the menu's open status and handles related actions.
     * If the menu is being closed, it emits a 'focusout' event and resets the input value.
     *
     * @param {boolean} value - The new open status of the menu. `true` to open, `false` to close.
     */
    const updateMenuOpened = (value) => {
      if (!value) {
        emit('focusout');
        setInputValue('');
      }

      opened.value = value;
    };

    /**
     * Sets the `opened` state to true, indicating that the menu should be displayed.
     */
    const showMenu = () => updateMenuOpened(true);

    /**
     * Focuses on the input element, bringing it into view for user interaction.
     */
    const focusInput = () => {
      inputElement.value?.focus();
      showMenu();
      emit('focus');
    };

    /**
     * Emits a 'click' event if the chip is not disabled.
     * This function checks the `disabled` prop before emitting the event.
     */
    const clickChip = () => !props.disabled && emit('click');

    /**
     * Removes a child chip from the selected items list based on its index.
     *
     * @param {number} index - The index of the chip to be removed from the selected items.
     */
    const closeChildChip = index => changeSelectedItems(
      selectedItems.value.filter((_, itemIndex) => itemIndex !== index),
    );

    /**
     * Emits a 'close' event to signal that the chip should be closed.
     * This function does not perform any checks before emitting the event.
     */
    const closeChip = () => emit('close');

    /**
     * LISTENERS FOR SINGLE CHIP
     */
    const chipListeners = computed(() => {
      const listeners = {};

      if (!props.disabled) {
        listeners.click = clickChip;
        listeners['click:close'] = closeChip;
      }

      return listeners;
    });

    /**
     * LISTENERS FOR MULTIPLE CHIPS
     */
    const multipleChipListeners = computed(() => {
      const listeners = {};

      if (!props.disabled) {
        listeners['click.prevent'] = () => {};
        listeners['click:close'] = closeChip;
      }

      return listeners;
    });

    /**
     * Selects an item and updates the selected items list based on the `multiple` prop.
     * Closes the menu and emits a 'next' event if `multiple` is false.
     * If `multiple` is true, it reopens the menu for additional selections.
     *
     * @param {*} value - The value of the item to be selected. It can be of any type depending on the context.
     */
    const selectItem = (value) => {
      opened.value = false;

      let newValue = value;

      if (props.multiple) {
        const prevValue = isArray(props.value) ? props.value : [props.value];
        const index = prevValue.findIndex(item => (
          (item?.[props.itemValue] && value?.[props.itemValue] && item?.[props.itemValue] === value?.[props.itemValue])
          || (item?.[props.itemValue] && item?.[props.itemValue] === value)
          || (value?.[props.itemValue] && value?.[props.itemValue] === item)
          || value === item
        ));

        newValue = [...prevValue];

        if (index === -1) {
          newValue.push(value);
        } else {
          newValue.splice(index, 1);
        }
      }

      changeSelectedItems(newValue);

      if (!props.multiple) {
        updateMenuOpened(false);
        emit('next');

        return;
      }

      updateSearch('');
      updateMenuOpened(false);
      nextTick(focusInput);
    };

    /**
     * Handles the keydown event for the input element, specifically responding to the Enter key.
     * If the Enter key is pressed, it either selects the current input value as an item or finds
     * and selects an item from the list that matches the input value.
     *
     * @param {KeyboardEvent} event - The keydown event triggered by the user.
     */
    const keydownInput = (event) => {
      if (event.keyCode === KEY_CODES.enter) { // TODO: change keyCode in whole application
        if (props.allowText) {
          const preparedValue = inputValue.value ?? '';

          selectItem(props.number ? Number(preparedValue) : preparedValue);
          return;
        }

        const lowerInputValue = inputValue.value.toLowerCase();
        const selectedItem = props.items.find(item => item[props.itemText]?.toLowerCase().startsWith(lowerInputValue));

        if (selectedItem) {
          selectItem(selectedItem);
        }
      }
    };

    /**
     * Checks if there are more items to fetch and triggers the fetch operation.
     * This function is used to load additional items when the user requests more data.
     * It relies on the `hasMoreItems` reactive property to determine if more items are available.
     */
    const showMore = () => {
      if (hasMoreItems.value) {
        fetchMoreItems();
      }
    };

    /**
     * Delays focusing on the input element by 100 milliseconds.
     * This function uses `setTimeout` to ensure the input element is focused after any pending operations.
     * It checks if the `inputElement` is available before attempting to focus.
     *
     * @todo Analise the solution with setTimeout in the future
     */
    const callFocus = () => setTimeout(() => (
      (props.active || props.alwaysActive) && focusInput()
    ), 100);

    watch(() => props.active, (active, prevActive) => {
      if (active && prevActive !== active) {
        setInputValue(selectedItems.value[0]?.[props.itemText] ?? selectedItems.value[0]?.[props.itemValue] ?? '', false);
        focusInput();
      }
    });

    /**
     * We need to fetch items when the items list changes after fetching (example: infos)
     */
    watch(() => props.items, (newItems, oldItems) => {
      const newItemsWithItemsValue = newItems.filter(item => item?.[props.childrenKey]?.length);
      const oldItemsWithItemsValue = oldItems.filter(item => item?.[props.childrenKey]?.length);

      if (newItemsWithItemsValue.length !== oldItemsWithItemsValue.length) {
        fetchItems();
      }
    });

    onMounted(() => {
      if (props.focusOnMount) {
        callFocus();
        focusRegister.register?.(callFocus);
      } else if (props.first) {
        focusRegister.register?.(callFocus);
      }
    });

    onBeforeUnmount(() => focusRegister.unregister?.(callFocus));

    return {
      inputPlaceholder,
      noDataText,
      hasAnyDisabledItem,
      selectedItems,
      opened,
      inputElement,
      inputValue,
      menuProps,
      inputType,
      chipText,
      lazyItems,
      valuesPending,
      pending,
      chipListeners,
      multipleChipListeners,
      hasNoDataFlagInSelectedItems,
      icon,

      selectItem,
      keydownInput,
      focusInput,
      callFocus,
      updateMenuOpened,
      updateInputValue,
      showMore,
      closeChildChip,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-chip.error input {
  color: var(--v-text-dark-primary, #FFFFFF);
}
</style>
