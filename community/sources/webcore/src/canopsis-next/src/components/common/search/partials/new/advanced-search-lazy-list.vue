<template>
  <v-list class="pa-0">
    <slot name="append" />
    <template v-for="item in items">
      <v-subheader
        v-if="item.header"
        :key="item.header"
      >
        {{ item.header }}
      </v-subheader>
      <v-list-item
        v-else
        :key="item.value"
        :input-value="isActiveItem(item)"
        @click="selectVariable(item)"
      >
        <v-list-item-content>
          <v-list-item-title>
            <v-layout class="gap-4" justify-space-between>
              <v-list-item-mask v-if="item[itemText]" :text="item[itemText]" :mask="search" />
              <span
                v-if="showValue"
                class="grey--text lighten-1"
              >
                {{ item[itemValue] }}
              </span>
            </v-layout>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
    <slot v-if="!items.length" name="no-data" />
    <v-menu
      v-if="subItemsShown"
      v-model="subItemsShown"
      :position-x="subItemsPosition.x"
      :position-y="subItemsPosition.y"
      :z-index="zIndex"
      offset-x
      right
    >
      <advanced-search-lazy-list
        :value="value"
        :items="parentItem[childrenKey]"
        :item-text="itemText"
        :item-value="itemValue"
        :z-index="zIndex + 1"
        :show-value="showValue"
        :children-key="childrenKey"
        @input="selectSubVariable"
      />
    </v-menu>
  </v-list>
</template>
<script>
import { uniqBy } from 'lodash';
import { ref } from 'vue';

export default {
  name: 'advanced-search-lazy-list', // TODO: see variables-list
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    search: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    childrenKey: {
      type: String,
      default: 'children',
    },
    selectedItems: {
      type: Array,
      default: () => [],
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    showValue: {
      type: Boolean,
      default: false,
    },
    zIndex: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const subItemsShown = ref(false);
    const parentItem = ref(undefined);
    const subItemsPosition = ref({ x: 0, y: 0 });

    const isActiveItem = (item = {}) => props.value.find((selectedItem) => {
      if (selectedItem[props.itemValue].length > String(item[props.itemValue]).length) {
        return selectedItem[props.itemValue].startsWith(`${item[props.itemValue]}.`);
      }

      return selectedItem[props.itemValue] === item[props.itemValue];
    });

    const selectVariable = item => emit('input', props.multiple ? uniqBy([...props.value, item], props.itemValue) : item); // TODO: refactor

    const selectSubVariable = (item) => {
      selectVariable(item);
      subItemsShown.value = false;
    };

    const handleMouseEnter = (item, event) => {
      if (item[this.childrenKey]) {
        const { left, top, width } = event.target.getBoundingClientRect();

        subItemsPosition.value.x = left + width;
        subItemsPosition.value.y = top;
        parentItem.value = item;
        subItemsShown.value = true;
      } else {
        parentItem.value = undefined;
        subItemsShown.value = false;
      }
    };

    return {
      subItemsShown,
      subItemsPosition,
      parentItem,

      handleMouseEnter,

      isActiveItem,
      selectVariable,
      selectSubVariable,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-subheader {
  font-size: 16px;
  font-weight: 700;
  color: inherit;
}
</style>
