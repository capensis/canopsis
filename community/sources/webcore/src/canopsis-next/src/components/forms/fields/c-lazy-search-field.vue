<template lang="pug">
  c-select-field.c-lazy-search-field(
    v-field="value",
    :search-input="search",
    :label="label",
    :loading="loading",
    :items="items",
    :name="name",
    :item-text="getItemText",
    :item-value="itemValue",
    :item-disabled="itemDisabled",
    :multiple="isMultiple",
    :deletable-chips="isMultiple",
    :small-chips="isMultiple",
    :chips="isMultiple",
    :disabled="disabled",
    :required="required",
    :menu-props="menuProps",
    :clearable="clearable",
    :autocomplete="autocomplete",
    :combobox="!autocomplete",
    :return-object="returnObject",
    :no-data-text="noDataText",
    no-filter,
    dense,
    @focus="onFocus",
    @blur="onBlur",
    @update:searchInput="updateSearch"
  )
    template(#item="{ item, tile, parent }")
      slot(name="item", :item="item", :tile="tile", :parent="parent")
        v-list-tile.c-lazy-search-field--tile(v-bind="tile.props", v-on="tile.on")
          slot(name="icon", :item="item")
          v-list-tile-content
            v-list-tile-mask(:text="getItemText(item)", :mask="internalSearch")
          span.ml-4.grey--text {{ item.type }}
    template(#append-item="")
      div.c-lazy-search-field__append(ref="append")
    template(v-if="$scopedSlots.selection", #selection="{ item, index }")
      slot(name="selection", :item="item", :index="index")
    template(v-else-if="isMultiple", #selection="{ item, index, parent }")
      v-chip.c-lazy-search-field__chip(
        v-if="isMultiple",
        small,
        close,
        @input="parent.onChipInput(item)"
      )
        span.ellipsis {{ getItemText(item) }}
</template>

<script>
import {
  debounce,
  isArray,
  isString,
  isFunction,
  get,
} from 'lodash';

export default {
  inject: ['$validator'],
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
  data() {
    return {
      isFocused: false,
      internalSearch: this.search,
    };
  },
  computed: {
    isMultiple() {
      return isArray(this.value);
    },

    menuProps() {
      return { contentClass: 'c-lazy-search-field__list' };
    },
  },
  watch: {
    search(value) {
      this.setInternalSearch(value);
    },

    isMultiple() {
      this.updateSearch(null, true);
    },
  },
  created() {
    this.debouncedEmitUpdateSearch = debounce(this.emitUpdateSearch, 300);
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

    setInternalSearch(value) {
      this.internalSearch = value;
    },

    emitUpdateSearch(value) {
      this.$emit('update:search', value);
    },

    updateSearch(value, force = false) {
      if (this.errors.has(this.name)) {
        this.errors.remove(this.name);
      }

      this.setInternalSearch(value);

      if (force) {
        this.emitUpdateSearch(value);
      } else if (this.isFocused) {
        this.debouncedEmitUpdateSearch(value);
      }
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

  &__chip {
    max-width: 100%;

    .v-chip__content {
      max-width: 100%;
    }
  }
}
</style>
