<template>
  <component
    :is="fieldComponent"
    :model-value="value"
    v-bind="fieldAttrs"
    :search="internalSearch"
    :label="label"
    :loading="loading"
    :items="items"
    :name="name"
    :item-title="getItemText"
    :item-value="getItemValue"
    :item-props="getItemProps"
    :multiple="isMultiple"
    :closable-chips="isMultiple"
    :chips="isMultiple"
    :disabled="disabled"
    :required="required"
    :menu-props="preparedMenuProps"
    :clearable="clearable"
    :return-object="returnObject"
    :no-data-text="noDataText"
    class="c-lazy-search-field"
    density="compact"
    no-filter
    @focus="onFocus"
    @blur="onBlur"
    @update:model-value="updateValue"
    @update:search="updateSearch"
  >
    <template #item="slotProps">
      <slot
        v-bind="slotProps"
        :attrs="slotProps.props"
        :item="slotProps.item?.raw ?? slotProps.item"
        name="item"
      >
        <v-list-item class="c-lazy-search-field--tile" v-bind="slotProps.props">
          <template v-if="$slots.icon" #prepend>
            <slot :item="slotProps.item?.raw ?? slotProps.item" name="icon" />
          </template>
          <template v-if="withType" #append>
            <span class="ml-4 text-medium-emphasis">
              {{ (slotProps.item?.raw ?? slotProps.item).type }}
            </span>
          </template>
        </v-list-item>
      </slot>
    </template>
    <template #append-item="">
      <div ref="appendElement" class="c-lazy-search-field__append" />
    </template>
    <template v-if="$slots.selection" #selection="slotProps">
      <slot name="selection" v-bind="slotProps" />
    </template>
    <template v-else-if="isMultiple" #selection="{ item, props: selectionProps }">
      <v-chip
v-bind="selectionProps"
class="c-lazy-search-field__chip"
size="small"
closable>
        <span class="text-truncate">
          {{ getItemText(getSelectionItem(item)) }}
        </span>
      </v-chip>
    </template>
    <template v-else #selection="{ item }">
      <span class="v-select__selection-text">
        {{ getItemText(getSelectionItem(item)) }}
      </span>
    </template>
    <template v-if="$slots['no-data']" #no-data="">
      <slot name="no-data" />
    </template>
  </component>
</template>

<script>
import { debounce, isArray, isString, isFunction, get } from 'lodash';
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue';

export default {
  inheritAttrs: false,
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
    itemTitle: {
      type: [String, Function],
      default: '',
    },
    itemValue: {
      type: [String, Function],
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
    withType: {
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
    menuProps: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ['update:value', 'update:modelValue', 'update:search', 'fetch', 'fetch:more'],
  setup(props, { emit, attrs }) {
    const internalSearch = ref(props.search);
    const isFocused = ref(false);
    const appendElement = ref(null);

    const fieldComponent = computed(() => (props.autocomplete ? 'c-autocomplete' : 'c-combobox'));
    const isMultiple = computed(() => isArray(props.value));
    const preparedMenuProps = computed(() => ({
      ...props.menuProps,
      contentClass: 'c-lazy-search-field__list',
      eager: true,
    }));
    const fieldAttrs = computed(() => attrs);

    const resolvedItemTitle = computed(() => props.itemTitle || props.itemText);

    const getItemValue = (item) => {
      if (isString(item)) {
        return item;
      }

      return isFunction(props.itemValue) ? props.itemValue(item) : get(item, props.itemValue);
    };

    const getItemText = (item) => {
      if (isString(item)) {
        return item;
      }

      const text = isFunction(resolvedItemTitle.value)
        ? resolvedItemTitle.value(item)
        : get(item, resolvedItemTitle.value);

      return text ?? getItemValue(item);
    };

    const getSelectionItem = (item) => {
      const rawItem = item?.raw ?? item;
      const value = getItemValue(rawItem);

      return props.items.find(entity => getItemValue(entity) === value) ?? rawItem;
    };

    const getItemProps = (item) => {
      const rawItem = item?.raw ?? item;

      if (isFunction(props.itemDisabled)) {
        return { disabled: props.itemDisabled(rawItem) };
      }

      return {};
    };

    const emitUpdateSearch = value => emit('update:search', value);
    const debouncedEmitUpdateSearch = debounce(emitUpdateSearch, props.debounce);

    const updateSearch = (value, force = false) => {
      internalSearch.value = value;

      if (force) {
        emitUpdateSearch(value);
      } else if (isFocused.value) {
        debouncedEmitUpdateSearch(value);
      }
    };

    const updateValue = (value) => {
      emit('update:value', value);
      emit('update:modelValue', value);
    };

    const onFocus = () => {
      isFocused.value = true;

      if (!props.items.length) {
        emit('fetch');
      }
    };

    const onBlur = () => {
      isFocused.value = false;
    };

    const intersectionHandler = (entries) => {
      const [entry] = entries;

      if (entry.isIntersecting && props.hasMore) {
        emit('fetch:more');
      }
    };

    const observer = new IntersectionObserver(intersectionHandler);

    watch(
      () => props.search,
      (value) => {
        internalSearch.value = value;
      },
    );

    watch(isMultiple, () => updateSearch(null, true));

    onMounted(() => {
      if (appendElement.value) {
        observer.observe(appendElement.value);
      }
    });
    onBeforeUnmount(() => {
      if (appendElement.value) {
        observer.unobserve(appendElement.value);
      }

      observer.disconnect();
    });

    return {
      appendElement,
      internalSearch,
      fieldComponent,
      fieldAttrs,
      preparedMenuProps,
      isMultiple,
      getItemValue,
      getItemText,
      getSelectionItem,
      getItemProps,
      updateSearch,
      updateValue,
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
    opacity: 0;
    height: 30px;
  }

  .v-select__selection {
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
