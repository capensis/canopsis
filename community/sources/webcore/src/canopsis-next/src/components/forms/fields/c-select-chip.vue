<template>
  <span class="c-select-chip d-inline-flex align-self-start">
    <v-menu :location="location" :disabled="disabled" offset="4">
      <template #activator="{ props: activatorProps }">
        <c-chip
          v-bind="{ ...chipAttrs, ...activatorProps }"
          :closable="!isEmptyValue && clearable"
          :text-color="chipTextColor"
          class="px-2"
          @close="selectItem(null)"
        >
          <slot name="selection">
            <slot v-if="isEmptyValue" name="selection-empty" />
            <span v-else>
              <slot name="selection-prefix" />
              {{ selectedItemText }}
              <c-help-icon
                v-if="selectedItemTooltip"
                :text="selectedItemTooltip"
                icon="help"
                icon-class="ml-1 grey--text"
                top
                small
              />
            </span>
          </slot>
        </c-chip>
      </template>
      <slot>
        <v-list density="compact">
          <v-list-item
            v-for="item in items"
            :key="item[itemValue]"
            :active="isSelectedItem(item)"
            :active-class="activeListItemClass"
            :disabled="item.disabled"
            @click="selectItem(returnObject ? item : item[itemValue])"
          >
            <template #title>
              <div>{{ item[itemText] }}</div>
              <div v-if="item.disabledMessage" class="text-caption">
                {{ item.disabledMessage }}
              </div>
            </template>
            <template v-if="item.tooltip" #append>
              <c-help-icon
                :text="item.tooltip"
                class="ml-2"
                icon="help"
                top
                small
              />
            </template>
          </v-list-item>
        </v-list>
      </slot>
    </v-menu>
  </span>
</template>

<script>
import { isNil, keyBy } from 'lodash';
import { computed, useAttrs } from 'vue';

export default {
  inheritAttrs: false,
  emits: ['update:value'],
  props: {
    value: {
      type: [String, Number, Object],
      required: false,
    },
    emptyLabel: {
      type: String,
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    top: {
      type: Boolean,
      default: false,
    },
    bottom: {
      type: Boolean,
      default: false,
    },
    left: {
      type: Boolean,
      default: false,
    },
    right: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    activeListItemClass: {
      type: String,
      default: '',
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const attrs = useAttrs();
    const isEmptyValue = computed(() => isNil(props.value));
    const itemsByValue = computed(() => keyBy(props.items, props.itemValue));
    const selectedValue = computed(() => props.value?.[props.itemValue] ?? props.value);
    const selectedItem = computed(() => itemsByValue.value[selectedValue.value]);
    const selectedItemText = computed(() => selectedItem.value?.[props.itemText] ?? selectedValue.value);
    const selectedItemTooltip = computed(() => selectedItem.value?.tooltip);
    const chipAttrs = computed(() => {
      const rest = { ...attrs };

      delete rest.textColor;
      delete rest['text-color'];

      return rest;
    });
    const chipTextColor = computed(() => {
      if (isEmptyValue.value) {
        return '';
      }

      return attrs.textColor ?? attrs['text-color'];
    });
    const location = computed(() => {
      const result = {
        [props.left]: 'left',
        [props.bottom]: 'bottom',
        [props.right]: 'right',
        [props.top]: 'top',
      }.true;

      return result ?? 'bottom';
    });

    const isSelectedItem = item => item[props.itemValue] === selectedValue.value;
    const selectItem = value => emit('update:value', value);

    return {
      isEmptyValue,
      selectedItemText,
      selectedItemTooltip,
      location,
      chipAttrs,
      chipTextColor,

      isSelectedItem,
      selectItem,
    };
  },
};
</script>
