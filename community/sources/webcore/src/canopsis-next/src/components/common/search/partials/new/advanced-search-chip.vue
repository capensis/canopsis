<template>
  <v-menu
    v-bind="menuProps"
    :value="opened"
    bottom
    @input="updateMenuOpened"
  >
    <template #activator="{ on }">
      <span v-if="multiple">
        <v-chip
          :close="closable && !disabled"
          :color="color"
          class="c-new-advanced-search__array-chip"
          v-on="multipleChipListeners"
        >
          <v-chip
            v-for="item in selectedItems"
            :key="item[itemValue]"
            :close="!disabled"
          >
            {{ item[itemText] ?? item[itemValue] }}
          </v-chip>
          <span>
            <input
              v-model="inputValue"
              ref="inputElement"
              type="text"
              class="ml-1"
              autocomplete="off"
              @keydown="keydown"
              @focus="focus"
              @input="updateInputValue"
            >
          </span>
        </v-chip>
      </span>
      <span v-else>
        <input
          v-show="active || alwaysActive"
          ref="inputElement"
          :value="inputValue"
          type="text"
          class="ml-1"
          autocomplete="off"
          v-on="on"
          @keydown="keydown"
          @focus="focus"
          @input="updateInputValue"
        >
        <v-chip
          v-if="!active && !alwaysActive"
          :color="color"
          :close="closable && !disabled"
          v-on="chipListeners"
        >
          <v-progress-circular
            v-if="valuesPending"
            color="primary"
            size="16"
            width="2"
            indeterminate
          />
          <span v-else>{{ chipText }}</span>
        </v-chip>
      </span>
    </template>
    <advanced-search-lazy-list
      :value="selectedItems"
      :items="lazyItems"
      :search="inputValue"
      :pending="fetchItems && pending"
      :item-text="itemText"
      :item-value="itemValue"
      @input="select"
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
    </advanced-search-lazy-list>
  </v-menu>
</template>

<script>
import {
  computed,
  ref,
  watch,
  toRef,
  inject,
  onMounted,
} from 'vue';

import { KEY_CODES } from '@/constants';

import { filterAdvancedSearchItems } from '@/helpers/search/advanced-search';

import { useI18n } from '@/hooks/i18n';
import { useLazySearch } from '@/hooks/form/lazy-search';

import AdvancedSearchLazyList from '@/components/common/search/partials/new/advanced-search-lazy-list.vue';

export default {
  components: { AdvancedSearchLazyList },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Number, Array],
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
  },
  setup(props, { emit }) {
    const registerLastInputFocus = inject('$registerLastInputFocus', () => {});

    const { t } = useI18n();

    const inputElement = ref(null);
    const inputValue = ref('');
    const opened = ref(false);

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

    const hasAnyDisabledItem = computed(() => props.items.some(({ disabled }) => disabled));

    const noDataText = computed(() => t(
      props.allowText
        ? 'common.pressEnterToApply'
        : 'common.noData',
    ));

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

    const {
      selectedItems,
      items: lazyItems,
      changeSelectedItems,
      updateSearch,
      valuesPending,
      wholePending: pending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
    } = useLazySearch({
      value: toRef(props, 'value'),
      addable: props.allowText,
      idKey: toRef(props, 'itemValue'),
      idParamsKey: 'ids',
      fetchHandler: props.fetchItems ?? defaultFetchItems,
      multiple: props.multiple,
    }, emit);

    const showMenu = () => opened.value = true;
    const focus = () => {
      showMenu();
      emit('focus');
    };

    const setInputValue = (newInputValue) => {
      inputValue.value = newInputValue;

      updateSearch(newInputValue);
    };

    const updateInputValue = event => setInputValue(event.target.value);

    const updateMenuOpened = (value) => {
      if (!value) {
        emit('focusout');
        setInputValue('');
      }

      opened.value = value;
    };

    const clickChip = () => !props.disabled && emit('click');
    const closeChip = () => emit('close');

    const chipText = computed(() => selectedItems.value
      .map(item => item[props.itemText] ?? item[props.itemValue])
      .join(','));

    const chipListeners = computed(() => {
      const listeners = {};

      if (!props.disabled) {
        listeners.click = clickChip;
        listeners['click:close'] = closeChip;
      }

      return listeners;
    });

    const multipleChipListeners = computed(() => {
      const listeners = {};

      if (!props.disabled) {
        listeners['click.prevent'] = () => {};
        listeners['click:close'] = closeChip;
      }

      return listeners;
    });

    const select = (value) => {
      opened.value = false;

      changeSelectedItems(props.multiple ? [...(props.value || []), value] : value);

      if (!props.multiple) {
        updateMenuOpened(false);
        emit('next');

        return;
      }

      updateMenuOpened(false);
      showMenu(); // TODO: think about it
    };

    const keydown = (event) => {
      if (event.keyCode === KEY_CODES.enter) {
        if (props.allowText) {
          select(inputValue.value ?? '');

          return;
        }

        const lowerInputValue = inputValue.value.toLowerCase();
        const selectedItem = props.items.find(item => item[props.itemText]?.toLowerCase().startsWith(lowerInputValue));

        if (selectedItem) {
          select(selectedItem);
        }
      }
    };

    const showMore = () => {
      if (hasMoreItems.value) {
        fetchMoreItems();
      }
    };

    const focusInput = () => setTimeout(() => inputElement.value?.focus(), 100); // TODO: check it

    watch(() => props.active, (active, prevActive) => {
      if (active && prevActive !== active) {
        setInputValue(selectedItems.value[0]?.[props.itemText] ?? selectedItems.value[0]?.[props.itemValue] ?? '');
        focusInput();
      }
    });

    /* watch(() => [props.active, selectedItems.value], ([active, newSelectedItems], [prevActive] = []) => {
      if (!active) {
        setInputValue('');

        return;
      }

      if (prevActive !== active) {
        nextTick(() => inputElement.value?.focus());
      }

      setInputValue(newSelectedItems[0]?.[props.itemText] ?? newSelectedItems[0]?.[props.itemValue] ?? '');
    }, { immediate: true }); */

    watch(() => props.items, fetchItems);

    onMounted(() => {
      if (!props.first) {
        setTimeout(() => inputElement.value?.focus(), 100);
      }

      registerLastInputFocus(focusInput);
    });

    return {
      noDataText,
      hasAnyDisabledItem,
      selectedItems,
      opened,
      inputElement,
      inputValue,
      menuProps,
      chipText,
      lazyItems,
      valuesPending,
      pending,
      chipListeners,
      multipleChipListeners,

      select,
      keydown,
      focus,
      updateMenuOpened,
      updateInputValue,
      showMore,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-chip.error input {
  color: var(--v-text-dark-primary, #FFFFFF);
}
</style>
