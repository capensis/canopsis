<template>
  <v-list
    :dense="dense"
    class="pa-0 variables-list"
  >
    <v-fade-transition>
      <v-progress-linear
        v-if="pending"
        class="variables-list__progress"
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
        :disabled="item.disabled || item[childrenKey]?.length === 0"
        @click="selectVariable(item)"
        @mouseenter="handleMouseEnter(item, $event)"
      >
        <v-list-item-content>
          <v-list-item-title>
            <v-layout class="gap-4" justify-space-between>
              <v-list-item-mask v-if="item[itemText]" :text="item[itemText]" :mask="search" />
              <span
                v-if="showValue && (!hideEmptyValue || !isUndefined(item[itemValue]))"
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
      class="variables-list__append-item"
    />
    <slot v-if="!items.length" name="no-data" />
    <v-menu
      v-if="subItemsShown"
      v-model="subItemsShown"
      :position-x="subItemsPosition.x"
      :position-y="subItemsPosition.y"
      :z-index="zIndex"
      :close-on-content-click="clickableParent"
      offset-x
      right
    >
      <variables-list
        :value="value"
        :items="parentItem[childrenKey]"
        :item-text="itemText"
        :item-value="itemValue"
        :z-index="zIndex + 1"
        :show-value="showValue"
        :hide-empty-value="hideEmptyValue"
        :children-key="childrenKey"
        :return-object="returnObject"
        :clickable-parent="clickableParent"
        @input="selectSubVariable"
      />
    </v-menu>
  </v-list>
</template>
<script>
import { uniq, uniqBy, isObject, isUndefined } from 'lodash';
import { ref, onMounted, onBeforeUnmount } from 'vue';

export default {
  name: 'variables-list', // We need it for recursive use
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String, Number],
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
    hideEmptyValue: {
      type: Boolean,
      default: false,
    },
    zIndex: {
      type: Number,
      required: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    clickableParent: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const appendElement = ref(null);
    const subItemsShown = ref(false);
    const parentItem = ref(undefined);
    const subItemsPosition = ref({ x: 0, y: 0 });

    /**
     * Prepares an item for use based on the component's configuration.
     *
     * @param {Object|string} item - The item to process, which can be an object or a string.
     * @returns {Object|string} The processed item. If `props.returnObject` is true, it returns the item itself;
     *                          otherwise, it returns the property specified by `props.itemValue` from the object,
     *                          or the item itself if it's not an object.
     */
    const prepareItem = (item) => {
      if (props.returnObject) {
        return item;
      }

      return isObject(item) ? item?.[props.itemValue] : item;
    };

    /**
     * Retrieves the value based on the component's configuration.
     *
     * @param {Object|string} value - The value to process, which can be an object or a string.
     * @returns {Object|string} The processed value. If `props.returnObject` is true, it returns the property
     *                          specified by `props.itemValue` from the object; otherwise, it returns the value itself.
     */
    const getValue = value => (props.returnObject ? value?.[props.itemValue] : value);

    /**
     * Determines if a given item is active based on the current selection.
     *
     * @param {Object} [item={}] - The item to check for activity status. Defaults to an empty object.
     * @returns {boolean} - Returns `true` if the item is active, otherwise `false`.
     */
    const isActiveItem = (item = {}) => props.value === item[props.itemValue] || props.value.find?.((selectedItem) => {
      const selectedValue = String(getValue(selectedItem) ?? '');
      const value = String(getValue(item) ?? '');

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
    const selectVariable = (item) => {
      if (item?.[props.childrenKey]?.length) {
        subItemsShown.value = true;
        return;
      }

      let newValue = prepareItem(item);

      if (props.multiple) {
        const newValueBeforeUniq = [...props.value, item];

        newValue = props.returnObject
          ? uniqBy(newValueBeforeUniq, props.itemValue)
          : uniq(newValueBeforeUniq);
      }

      emit('input', newValue);
    };

    /**
     * Selects a sub-variable by modifying its value to include the parent item's value
     * and then calls the `selectVariable` function with the modified item.
     *
     * @param {Object} item - The sub-item to be selected. This object should contain
     *                        properties that can be identified by the `props.itemValue`.
     */
    const selectSubVariable = (item) => {
      /**
       * We are replacing /\s*}}/ for correct parent prefix checking for payloads.
       * If parent value is not set, we are not adding parent prefix.
       * @example detect {{ .LastChild.Alarm.Value }} in {{ .LastChild.Alarm.Value.Infos.something.Value }}
       */
      const hasParentValue = !!parentItem.value[props.itemValue];
      const parentValue = String(parentItem.value[props.itemValue]).replace(/\s*}}$/, '');

      const value = getValue(item);
      const newValue = String(value).startsWith(parentValue) || !hasParentValue
        ? value
        : `${parentValue}.${value}`;

      if (props.returnObject) {
        selectVariable({
          ...item,
          [props.itemValue]: newValue,
        });
      } else {
        selectVariable(newValue);
      }

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

      isUndefined,
    };
  },
};
</script>

<style lang="scss" scoped>
.variables-list {
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
