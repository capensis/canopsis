<template>
  <v-list class="pa-0 advanced-search-lazy-list">
    <v-fade-transition>
      <v-progress-linear
        v-if="pending"
        class="advanced-search-lazy-list__progress"
        color="primary"
        indeterminate
      />
    </v-fade-transition>
    <slot :items="items" name="prepend" />
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
        :disabled="item.disabled"
        @click="selectVariable(item)"
        @mouseenter="handleMouseEnter(item, $event)"
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
        <v-list-item-action v-if="item[childrenKey]">
          <v-icon>arrow_right</v-icon>
        </v-list-item-action>
      </v-list-item>
    </template>
    <div
      ref="appendElement"
      class="advanced-search-lazy-list__append-item"
    />
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
import { ref, onMounted, onBeforeUnmount } from 'vue';

export default {
  name: 'advanced-search-lazy-list', // We need it for recursive use
  model: {
    prop: 'value',
    event: 'input',
  },
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
      default: 'items',
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
    const appendElement = ref(null);
    const subItemsShown = ref(false);
    const parentItem = ref(undefined);
    const subItemsPosition = ref({ x: 0, y: 0 });

    /**
     * Determines if a given item is active based on the current selection.
     *
     * @param {Object} [item={}] - The item to check for activity status. Defaults to an empty object.
     * @returns {boolean} - Returns `true` if the item is active, otherwise `false`.
     */
    const isActiveItem = (item = {}) => props.value.find((selectedItem) => {
      const selectedValue = String(selectedItem[props.itemValue] ?? '');
      const value = String(item[[props.itemValue]] ?? '');

      if (selectedValue.length > value.length) {
        return selectedValue.startsWith(`${value}.`);
      }

      return selectedValue === value;
    });

    /**
     * Emits an 'input' event with the selected item(s) based on the component's configuration.
     *
     * @param {Object} item - The item to be selected and emitted. This object should contain properties
     *                        that can be identified by the `props.itemValue`.
     */
    const selectVariable = item => emit(
      'input',
      props.multiple ? uniqBy([...props.value, item], props.itemValue) : item,
    );

    /**
     * Selects a sub-variable by modifying its value to include the parent item's value
     * and then calls the `selectVariable` function with the modified item.
     *
     * @param {Object} item - The sub-item to be selected. This object should contain
     *                        properties that can be identified by the `props.itemValue`.
     */
    const selectSubVariable = (item) => {
      selectVariable({
        ...item,
        [props.itemValue]: `${parentItem.value[props.itemValue]}.${item[props.itemValue]}`,
      });
      subItemsShown.value = false;
    };

    /**
     * Handles the mouse enter event on an item, updating the position and visibility
     * of sub-items if applicable.
     *
     * @param {Object} item - The item that the mouse has entered. This object should
     *                        contain properties that can be identified by the `props.childrenKey`.
     * @param {MouseEvent} event - The mouse event triggered by entering the item.
     */
    const handleMouseEnter = (item, event) => {
      if (item.disabled) {
        return;
      }

      if (item[props.childrenKey]) {
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

    /**
     * Handles intersection events to determine when more data should be fetched.
     *
     * @param {IntersectionObserverEntry[]} entries - An array of IntersectionObserverEntry objects,
     *                                                containing information about the intersection
     *                                                changes for the observed target elements.
     */
    const intersectionHandler = (entries) => {
      const [entry] = entries;

      if (entry.isIntersecting) {
        emit('fetch:more');
      }
    };

    const observer = new IntersectionObserver(intersectionHandler);

    onMounted(() => observer.observe(appendElement.value));
    onBeforeUnmount(() => {
      observer.unobserve(appendElement.value);
      observer.disconnect();
    });

    return {
      appendElement,
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
.advanced-search-lazy-list {
  position: relative;

  .v-subheader {
    font-size: 16px;
    font-weight: 700;
    color: inherit;
  }

  &__progress {
    position: sticky;
    top: 0;
    left: 0;
    width: 100%;
  }

  &__append-item {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    height: 48px;
    pointer-events: none;
  }
}
</style>
