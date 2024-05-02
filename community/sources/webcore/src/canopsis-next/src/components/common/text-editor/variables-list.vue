<template>
  <v-list
    :dense="dense"
    class="pa-0"
  >
    <slot name="append" />
    <v-list-item
      v-for="item in items"
      :key="item.value"
      :input-value="isActiveVariable(item)"
      @click="selectVariable(item)"
      @mouseenter="handleMouseEnter(item, $event)"
    >
      <v-list-item-content>
        <v-list-item-title>
          <v-layout class="gap-4" justify-space-between>
            <span v-if="item.text">{{ item.text }}</span>
            <span
              v-if="showValue"
              class="grey--text lighten-1"
            >
              {{ item.value }}
            </span>
          </v-layout>
        </v-list-item-title>
      </v-list-item-content>
      <v-list-item-action v-if="item[childrenKey]">
        <v-icon>arrow_right</v-icon>
      </v-list-item-action>
    </v-list-item>
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
      <variables-list
        :value="value"
        :items="parentItem[childrenKey]"
        :z-index="zIndex + 1"
        :show-value="showValue"
        :return-object="returnObject"
        :children-key="childrenKey"
        @input="selectSubVariable"
      />
    </v-menu>
  </v-list>
</template>

<script>
export default {
  name: 'VariablesList',
  props: {
    value: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    zIndex: {
      type: Number,
      required: false,
    },
    showValue: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    /**
     * TODO: rename `variables` in children to `items` in all components and remove this property in the future.
     */
    childrenKey: {
      type: String,
      default: 'variables',
    },
  },
  data() {
    return {
      subItemsShown: false,
      parentItem: undefined,
      subItemsPosition: {
        x: 0,
        y: 0,
      },
    };
  },
  methods: {
    selectVariable(variable) {
      this.$emit('input', this.returnObject ? variable : variable.value);
    },

    selectSubVariable(value) {
      this.$emit('input', value);
      this.subItemsShown = false;
    },

    isActiveVariable(item) {
      if (this.value.length > item.value.length) {
        return this.value.startsWith(`${item.value}.`);
      }

      return this.value.startsWith(item.value);
    },

    handleMouseEnter(item, event) {
      if (item[this.childrenKey]) {
        const { left, top, width } = event.target.getBoundingClientRect();

        this.subItemsPosition.x = left + width;
        this.subItemsPosition.y = top;
        this.parentItem = item;
        this.subItemsShown = true;
      } else {
        this.parentItem = undefined;
        this.subItemsShown = false;
      }
    },
  },
};
</script>
