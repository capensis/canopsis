<template>
  <div
    class="grid-layout"
    ref="layout"
    :class="{ 'grid-layout--disabled': disabled }"
    :style="style"
  >
    <grid-item
      class="grid-layout__placeholder primary darken-1"
      v-if="!disabled"
      v-show="resizing || moving"
      v-bind="itemBind"
      :item="placeholder"
      key="placeholder"
    />
    <grid-item
      v-for="layoutItem in layout"
      v-on="itemOn"
      v-bind="itemBind"
      :key="layoutItem.i"
      :item="layoutItem"
    >
      <template #default="props">
        <slot
          name="item"
          v-bind="props"
        />
      </template>
    </grid-item>
  </div>
</template>

<script>
import {
  compactLayout,
  calculateLayoutRowsCount,
  calculateLayoutBottom,
  findLayoutItem,
  moveLayoutItem,
  replaceLayoutItemInLayout,
} from '@/helpers/grid';

import { formBaseMixin } from '@/mixins/form';

import GridItem from './partials/grid-item.vue';

export default {
  components: { GridItem },
  mixins: [formBaseMixin],
  model: {
    prop: 'layout',
    event: 'input',
  },
  props: {
    layout: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    columnsCount: {
      type: Number,
      default: 12,
    },
    rowHeight: {
      type: Number,
      default: 150,
    },
    margin: {
      type: Array,
      default: () => [10, 10],
    },
    throttle: {
      type: Number,
      default: 10,
    },
    debounce: {
      type: Number,
      default: 50,
    },
  },
  data() {
    return {
      width: null,
      resizing: false,
      moving: false,
      placeholder: {
        x: 0,
        y: 0,
        w: 0,
        h: 0,
        i: '-1',
      },
    };
  },
  computed: {
    itemBind() {
      return {
        layout: this.layout,
        containerWidth: this.width,
        rowHeight: this.rowHeight,
        columnsCount: this.columnsCount,
        margin: this.margin,
        disabled: this.disabled,
        throttle: this.throttle,
        debounce: this.debounce,
      };
    },

    itemOn() {
      if (this.disabled) {
        return {};
      }

      return {
        resize: this.resizeItemHandler,
        resized: this.resizedItemHandler,
        move: this.moveItemHandler,
        moved: this.movedItemHandler,
        toggle: this.toggleAutoHeightHandler,
      };
    },

    disabledLayoutStyle() {
      const rowHeightInPixels = `${this.rowHeight}px`;

      return {
        padding: rowHeightInPixels,
        columnGap: rowHeightInPixels,
        gridTemplateColumns: `repeat(${this.columnsCount}, 1fr)`,
      };
    },

    enabledLayoutStyle() {
      const [, marginY] = this.margin;
      const rowCount = calculateLayoutRowsCount(this.layout);
      const bottom = calculateLayoutBottom(this.layout);

      return {
        height: `${(bottom * this.rowHeight) + (marginY * rowCount) + marginY}px`,
      };
    },

    style() {
      return this.disabled ? this.disabledLayoutStyle : this.enabledLayoutStyle;
    },
  },
  watch: {
    disabled(disabled) {
      if (disabled) {
        this.removeAllListeners();

        return;
      }

      this.addAllListeners();
      this.resizeObserverHandler();
    },
  },
  mounted() {
    if (!this.disabled) {
      this.addAllListeners();
      this.resizeObserverHandler();
    }
  },
  beforeDestroy() {
    this.removeAllListeners();
  },
  methods: {
    toggleAutoHeightHandler(layoutItem) {
      const newLayoutItem = { ...layoutItem, autoHeight: !layoutItem.autoHeight };
      const newLayout = compactLayout(replaceLayoutItemInLayout(this.layout, newLayoutItem));

      this.updateModel(newLayout);
    },

    /**
     * Resize grid item handler
     *
     * @param {GridLayoutItem} layoutItem
     * @param {number} w
     * @param {number} h
     */
    resizeItemHandler(layoutItem, w, h) {
      this.resizing = true;

      const newLayoutItem = { ...layoutItem, h, w };
      const newLayout = compactLayout(replaceLayoutItemInLayout(this.layout, newLayoutItem));

      this.placeholder = {
        ...this.placeholder,
        ...newLayoutItem,
      };

      this.updateModel(newLayout);
    },

    /**
     * Resize END grid item handler
     *
     * @param {GridLayoutItem} layoutItem
     * @param {number} w
     * @param {number} h
     */
    resizedItemHandler(layoutItem, w, h) {
      const newLayoutItem = { ...layoutItem, w, h };
      const newLayout = compactLayout(replaceLayoutItemInLayout(this.layout, newLayoutItem));

      this.updateModel(newLayout);
      this.resizing = false;
    },

    /**
     * Move grid item handler
     *
     * @param {GridLayoutItem} layoutItem
     * @param {number} x
     * @param {number} y
     */
    moveItemHandler(layoutItem, x, y) {
      this.moving = true;

      const newLayout = compactLayout(moveLayoutItem(this.layout, layoutItem, x, y, true));

      this.placeholder = {
        ...this.placeholder,
        ...findLayoutItem(newLayout, layoutItem.i),
      };

      this.updateModel(newLayout);
    },

    /**
     * Move END grid item handler
     *
     * @param {GridLayoutItem} layoutItem
     * @param {number} x
     * @param {number} y
     */
    movedItemHandler(layoutItem, x, y) {
      const newLayout = compactLayout(moveLayoutItem(this.layout, layoutItem, x, y, true));

      this.updateModel(newLayout);
      this.moving = false;
    },

    /**
     * Element resize observer handler
     */
    resizeObserverHandler() {
      const newWidth = this.$refs?.layout?.offsetWidth ?? null;

      if (newWidth !== this.width) {
        this.width = newWidth;
      }
    },

    /**
     * Add all component listeners
     */
    addAllListeners() {
      this.$resizeObserver = new ResizeObserver(this.resizeObserverHandler);
      this.$resizeObserver.observe(this.$el);
    },

    /**
     * Remove all component listeners
     */
    removeAllListeners() {
      this.$resizeObserver?.disconnect();
    },
  },
};
</script>

<style lang="scss" scoped>
.grid-layout {
  &:not(&--disabled) {
    background-color: rgba(60, 60, 60, 0.05);
    margin: auto auto 500px auto;
    position: relative;
  }

  &__placeholder {
    opacity: .35;

    &:after {
      content: none !important;
    }
  }

  &--disabled {
    display: grid;
  }
}
</style>
