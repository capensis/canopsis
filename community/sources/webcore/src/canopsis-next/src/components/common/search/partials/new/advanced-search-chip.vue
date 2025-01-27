<template>
  <v-menu
    v-bind="menuProps"
    v-model="opened"
    bottom
    @input="toggleMenu"
  >
    <template #activator="{ on }">
      <span v-if="multiple">
        <v-chip
          :close="closable"
          :color="color"
          class="c-new-advanced-search__array-chip"
          @click.prevent=""
          @click:close="close"
        >
          <v-chip
            v-for="item in selectedItems"
            :key="item[itemValue]"
            close
          >
            {{ item[itemText] ?? item[itemValue] }}
          </v-chip>
          <span>
            <input
              v-model="inputValue"
              ref="inputEl"
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
          ref="inputEl"
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
          :close="closable"
          @click="click"
          @click:close="close"
        >
          {{ chipText }}
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
      <template v-if="allowText" #no-data="">
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title v-html="$t('common.pressEnterToApply')" />
          </v-list-item-content>
        </v-list-item>
      </template>
    </advanced-search-lazy-list>
  </v-menu>
</template>

<script>
import { isUndefined } from 'lodash';
import {
  computed,
  ref,
  watch,
  toRef,
  nextTick,
  onMounted,
} from 'vue';

import { KEY_CODES } from '@/constants';

import { useLazySearch } from '@/hooks/form/lazy-search';

import AdvancedSearchLazyList from '@/components/common/search/partials/new/advanced-search-lazy-list.vue';

export const filterItems = (items, condition) => {
  let lastHeaderIndex;

  return items.reduce((acc, item, index) => {
    if (item.header) {
      lastHeaderIndex = index;

      return acc;
    }

    if (condition(item)) {
      if (!isUndefined(lastHeaderIndex)) {
        acc.push(items[lastHeaderIndex]);
      }

      acc.push(item);
    }

    return acc;
  }, []);
};

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
  },
  setup(props, { emit }) {
    const inputEl = ref(null);
    const appendEl = ref(null);
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
    }));

    const showMenu = () => opened.value = true;
    const hideMenu = () => opened.value = false;

    const apply = (event) => {
      console.log('APPLY:', event);
    };

    const defaultFetchItems = ({ params = {} } = {}) => {
      let data = props.items;

      if (params?.search) {
        const lowerSearch = String(params.search).toLowerCase();

        data = filterItems(data, item => item[props.itemText]?.toLowerCase().includes(lowerSearch));
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
      wholePending: pending,
      hasMoreItems,
      fetchMoreItems,
    } = useLazySearch({
      value: toRef(props, 'value'),
      addable: props.allowText,
      idKey: toRef(props, 'itemValue'),
      idParamsKey: 'ids',
      fetchHandler: props.fetchItems ?? defaultFetchItems,
      multiple: props.multiple,
    }, emit);

    const chipText = computed(() => selectedItems.value.map(item => item[props.itemText] ?? item[props.itemValue]).join(','));

    const click = () => emit('click');
    const close = () => emit('close');
    const focus = () => {
      showMenu();
      emit('focus');
    };
    // const focusout = () => emit('focusout');

    const updateInputValue = (event) => {
      inputValue.value = event.target.value;

      updateSearch(inputValue.value);
    };

    const toggleMenu = (value) => {
      if (!value) {
        emit('focusout');
        inputValue.value = '';
      }
    };

    const select = (value) => {
      opened.value = false;

      changeSelectedItems(props.multiple ? [...(props.value || []), value] : value);

      if (!props.multiple) {
        toggleMenu();
        emit('next');
      } else {
        toggleMenu();
        showMenu(); // TODO: think about it
      }
    };

    const keydown = (event) => {
      if (event.keyCode === KEY_CODES.enter && props.allowText) { // TODO: try to find in items
        select(inputValue.value ?? '');

        return;
      }
      console.log('KEYDOWN:', event);
    };

    watch(() => [props.active, selectedItems.value], ([active, newSelectedItems], [prevActive] = []) => {
      if (!active) {
        inputValue.value = '';
        updateSearch(inputValue.value);

        return;
      }

      console.log(active, prevActive, props.items, inputEl.value);

      if (prevActive !== active) {
        nextTick(() => inputEl.value?.focus());
      }

      inputValue.value = newSelectedItems[0]?.[props.itemText] ?? newSelectedItems[0]?.[props.itemValue] ?? '';
    }, { immediate: true });

    const showMore = () => {
      if (hasMoreItems.value) {
        fetchMoreItems();
      }
    };

    onMounted(() => {
      if (!props.first) {
        setTimeout(() => inputEl.value?.focus(), 100);
      }
    });

    return {
      selectedItems,
      opened,
      appendEl,
      inputEl,
      inputValue,
      menuProps,
      chipText,
      lazyItems,
      pending,

      apply,
      select,
      keydown,
      click,
      focus,
      toggleMenu,
      close,
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
