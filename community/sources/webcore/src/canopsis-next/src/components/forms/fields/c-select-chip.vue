<template>
  <v-menu
    :top="top"
    :bottom="bottom"
    :left="left"
    :right="right"
    offset-y
  >
    <template #activator="{ on }">
      <c-chip v-bind="$attrs" class="px-2" v-on="on">
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
      <v-list>
        <v-list-item
          v-for="item in items"
          :key="item[itemValue]"
          :disabled="item.disabled"
          @click="selectItem(item[itemValue])"
        >
          <v-list-item-content>
            <div>{{ item[itemText] }}</div>
            <div v-if="item.disabledText" class="text-caption">
              {{ item.disabledText }}
            </div>
          </v-list-item-content>
          <v-list-item-action>
            <c-help-icon
              v-if="item.tooltip"
              :text="item.tooltip"
              class="ml-2"
              icon="help"
              top
              small
            />
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </slot>
  </v-menu>
</template>

<script>
import { isNil, keyBy } from 'lodash';
import { computed } from 'vue';

export default {
  inheritAttrs: false,
  props: {
    value: {
      type: [String, Number],
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
  },
  setup(props, { emit }) {
    const isEmptyValue = computed(() => isNil(props.value));
    const itemsByValue = computed(() => keyBy(props.items, props.itemValue));
    const selectedItemText = computed(() => itemsByValue.value[props.value]?.[props.itemText] ?? props.value);
    const selectedItemTooltip = computed(() => itemsByValue.value[props.value]?.tooltip);

    const selectItem = value => emit('input', value);

    return {
      isEmptyValue,
      selectedItemText,
      selectedItemTooltip,

      selectItem,
    };
  },
};
</script>
