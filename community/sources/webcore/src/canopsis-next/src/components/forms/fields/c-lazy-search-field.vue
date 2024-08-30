<template>
  <c-select-field
    v-field="value"
    v-bind="$attrs"
    :search-input="search"
    :label="label"
    :loading="loading"
    :items="items"
    :name="name"
    :item-text="getItemText"
    :item-value="itemValue"
    :item-disabled="itemDisabled"
    :multiple="isMultiple"
    :deletable-chips="isMultiple"
    :small-chips="isMultiple"
    :chips="isMultiple"
    :disabled="disabled"
    :required="required"
    :menu-props="menuProps"
    :clearable="clearable"
    :autocomplete="autocomplete"
    :combobox="!autocomplete"
    :return-object="returnObject"
    :no-data-text="noDataText"
    class="c-lazy-search-field mt-4"
    no-filter
    dense
    @focus="onFocus"
    @blur="onBlur"
    @update:search-input="updateSearch"
  >
    <template #item="{ item, attrs, on, parent }">
      <slot
        :attrs="attrs"
        :on="on"
        :item="item"
        :parent="parent"
        name="item"
      >
        <v-list-item
          class="c-lazy-search-field--tile"
          v-bind="attrs"
          v-on="on"
        >
          <slot
            :item="item"
            name="icon"
          />
          <v-list-item-content>
            <v-list-item-mask :text="getItemText(item)" :mask="internalSearch" />
          </v-list-item-content>
          <span class="ml-4 grey--text">{{ item.type }}</span>
        </v-list-item>
      </slot>
    </template>
    <template #append-item="">
      <div
        ref="appendElement"
        class="c-lazy-search-field__append"
      />
    </template>
    <template v-if="$scopedSlots.selection" #selection="{ item, index, parent }">
      <slot
        :item="item"
        :index="index"
        :parent="parent"
        name="selection"
      />
    </template>
    <template v-else-if="isMultiple" #selection="{ item, parent }">
      <v-chip
        v-if="isMultiple"
        class="c-lazy-search-field__chip"
        small
        close
        @click:close="parent.onChipInput(item)"
      >
        <span class="text-truncate">{{ getItemText(item) }}</span>
      </v-chip>
    </template>
  </c-select-field>
</template>

<script>
import {
  debounce,
  isArray,
  isString,
  isFunction,
  get,
} from 'lodash';
import {
  computed,
  ref,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { useValidator } from '@/hooks/validator/validator';

export default {
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String, Object],
      default: '',
    },
    name: {
      type: String,
      default: 'entities',
    },
    label: {
      type: String,
      required: false,
    },
    itemText: {
      type: [String, Function],
      default: '_id',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    noDataText: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    autocomplete: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    itemDisabled: {
      type: [String, Array, Function],
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    hasMore: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    search: {
      type: String,
      default: null,
    },
    debounce: {
      type: Number,
      default: 300,
    },
  },
  setup(props, { emit }) {
    const internalSearch = ref(props.search);
    const isFocused = ref(false);
    const appendElement = ref(null);

    const validator = useValidator();

    const isMultiple = computed(() => isArray(props.value));
    const menuProps = computed(() => ({ contentClass: 'c-lazy-search-field__list', eager: true }));

    /**
     * Set the internal search value.
     * @param {any} value - The value to set as the internal search.
     */
    const setInternalSearch = value => internalSearch.value = value;

    /**
     * Emit an update search event.
     *
     * @param {any} value - The value to emit for the update search event.
     */
    const emitUpdateSearch = value => emit('update:search', value);

    /**
     * Debounced version of `emitUpdateSearch` based on the specified throttle time.
     *
     * @param {any} value - The value to emit for the update search event.
     */
    const debouncedEmitUpdateSearch = debounce(emitUpdateSearch, props.throttle);

    /**
     * Get the text of an item based on the provided item.
     *
     * @param {any} item - The item to get the text from.
     * @returns {string} The text of the item.
     */
    const getItemText = (item) => {
      if (isString(item)) {
        return item;
      }

      return isFunction(props.itemText) ? props.itemText(item) : get(item, props.itemText);
    };

    /**
     * Update the search value and emit the update event if conditions are met.
     *
     * @param {any} value - The new value for the search.
     * @param {boolean} [force = false] - Whether to force emit the update event.
     */
    const updateSearch = (value, force = false) => {
      if (validator.errors.has(props.name)) {
        validator.errors.remove(props.name);
      }

      setInternalSearch(value);

      if (force) {
        emitUpdateSearch(value);
      } else if (isFocused.value) {
        debouncedEmitUpdateSearch(value);
      }
    };

    /**
     * Handle the focus event.
     */
    const onFocus = () => {
      isFocused.value = true;

      if (!props.items.length) {
        emit('fetch');
      }
    };

    /**
     * Handle the blur event.
     */
    const onBlur = () => isFocused.value = false;

    /**
     * Handle the intersection observer callback for lazy loading more items.
     *
     * @param {IntersectionObserverEntry[]} entries - The entries observed by the intersection observer.
     */
    const intersectionHandler = (entries) => {
      const [entry] = entries;

      if (entry.isIntersecting && props.hasMore) {
        emit('fetch:more');
      }
    };

    const observer = new IntersectionObserver(intersectionHandler);

    watch(() => props.search, setInternalSearch);
    watch(isMultiple, () => updateSearch(null, true));

    onMounted(() => observer.observe(appendElement.value));
    onBeforeUnmount(() => observer.unobserve(appendElement.value));

    return {
      appendElement,

      internalSearch,
      isFocused,
      menuProps,

      isMultiple,

      getItemText,
      updateSearch,
      onFocus,
      onBlur,
    };
  },
};
</script>

<style lang="scss">
.c-lazy-search-field {
  &__list .v-list {
    position: relative;
  }

  &__append {
    position: absolute;
    pointer-events: none;
    right: 0;
    bottom: 0;
    left: 0;
    height: 200px;
  }

  .v-select__selections {
    max-width: calc(100% - 24px);
  }

  &.v-autocomplete:not(.v-input--is-focused).v-select--chips input {
    max-height: 32px;
  }

  .v-chip--select {
    max-width: 100%;

    .v-chip__content {
      max-width: 100%;
      white-space: nowrap !important;
      overflow: hidden !important;
      text-overflow: ellipsis !important;
    }
  }
}
</style>
