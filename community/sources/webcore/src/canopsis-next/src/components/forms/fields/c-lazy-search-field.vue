<template>
  <c-select-field
    v-field="value"
    v-validate="rules"
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
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    :menu-props="{ contentClass: 'c-lazy-search-field__list', eager: true }"
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
    @update:search-input="debouncedUpdateSearch"
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
            {{ getItemText(item) }}
          </v-list-item-content>
          <span class="ml-4 grey--text">{{ item.type }}</span>
        </v-list-item>
      </slot>
    </template>
    <template #append-item="">
      <div
        ref="append"
        class="c-lazy-search-field__append"
      />
    </template>
    <template #selection="{ item, index, parent }">
      <slot
        :item="item"
        :index="index"
        :parent="parent"
        name="selection"
      >
        <v-chip
          v-if="isMultiple"
          class="c-lazy-search-field__chip"
          small
          close
          @click:close="removeItemFromArray(index)"
        >
          <span class="text-truncate">{{ getItemText(item) }}</span>
        </v-chip>
        <slot
          v-else
          :item="item"
          name="selection"
        >
          {{ getItemText(item) }}
        </slot>
      </slot>
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

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
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
  },
  data() {
    return {
      isFocused: false,
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    isMultiple() {
      return isArray(this.value);
    },
  },
  created() {
    this.debouncedUpdateSearch = debounce(this.updateSearch, 300);
  },
  mounted() {
    this.observer = new IntersectionObserver(this.intersectionHandler);

    this.observer.observe(this.$refs.append);
  },
  beforeDestroy() {
    this.observer.unobserve(this.$refs.append);
  },
  methods: {
    getItemText(item) {
      if (isString(item)) {
        return item;
      }

      return isFunction(this.itemText) ? this.itemText(item) : get(item, this.itemText);
    },

    intersectionHandler(entries) {
      const [entry] = entries;

      if (entry.isIntersecting && this.hasMore) {
        this.fetchItems();
      }
    },

    updateSearch(value) {
      if (this.isFocused) {
        this.$emit('update:search', value);
      }
    },

    onFocus() {
      this.isFocused = true;

      if (!this.items.length) {
        this.$emit('fetch');
      }
    },

    onBlur() {
      this.isFocused = false;
    },

    fetchItems() {
      this.$emit('fetch:more');
    },

    updateValue(value) {
      this.updateModel(this.returnObject ? value : get(value, this.itemValue));
    },
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

  &__chip {
    max-width: 100%;

    .v-chip__content {
      max-width: 100%;
    }
  }
}
</style>
