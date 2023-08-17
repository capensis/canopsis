<template lang="pug">
  div.c-grid-layout(
    ref="layout",
    :style="preparedStyle",
    :class="{ 'c-grid-layout--disabled': disabled }"
  )
    c-grid-item.c-grid-layout__placeholder(
      v-if="!disabled",
      v-show="resizing || moving",
      v-bind="bind",
      :i="placeholder.i",
      :x="placeholder.x",
      :y="placeholder.y",
      :w="placeholder.w",
      :h="placeholder.h"
    )
    slot(:bind="bind", :on="on")
</template>

<script>
import { uniq } from 'lodash';

import { WIDGET_GRID_COLUMNS_COUNT, WIDGET_GRID_ROW_HEIGHT } from '@/constants';

import { formArrayMixin } from '@/mixins/form';

import CGridItem from './c-grid-item.vue';

/**
 * TODO: MOVE TO HELPERS
 */
const calculateLayoutBottom = (layout = []) => layout.reduce((acc, item) => Math.max(acc, item.y + item.h), 0);
const calculateLayoutRowCount = (layout = []) => uniq(layout.map(({ y }) => y)).length;

/**
 * @typedef GridLayoutItem
 * @property {number | string} i
 * @property {number} x
 * @property {number} y
 * @property {number} w
 * @property {number} h
 * @property {boolean} moved
 * @property {boolean} resized
 */

/**
 * @typedef {GridLayoutItem[]} GridLayout
 */

/**
 * Given two layout items, check if they collide.
 *
 * @param {GridLayoutItem} first
 * @param {GridLayoutItem} second
 * @returns {boolean}
 */
const collides = (first, second) => !(
  first.i === second.i // same element
  || first.x + first.w <= second.x // first is left of second
  || first.x >= second.x + second.w // first is right of second
  || first.y + first.h <= second.y // first is above second
  || first.y >= second.y + second.h // first is below second
);

/**
 * Get all collisions for layout
 *
 * @param {GridLayout} layout
 * @param {GridLayoutItem} layoutItem
 * @returns {GridLayout}
 */
const getAllCollisions = (layout, layoutItem) => layout.filter(l => collides(layoutItem, l));

/**
 * Get sorted layout by rows and columns
 *
 * @param {GridLayout} layout
 * @returns {GridLayout}
 */
const getSortedLayout = (layout = []) => (
  [...layout].sort((a, b) => Number(a.y > b.y || (a.y === b.y && a.x > b.x)) || -1)
);

const getFirstCollision = (layout, layoutItem) => layout.find(item => collides(item, layoutItem));

/**
 *
 * @param {GridLayout} layout
 * @param {GridLayoutItem} layoutItem
 * @param {number} x
 * @param {number} y
 * @param {boolean} isUserAction
 * @returns {*}
 */
function moveElement(layout, layoutItem, x, y, isUserAction) {
  const newLayoutItem = {
    ...layoutItem,
    x,
    y,
    moved: true,
  };
  const newLayout = layout.map(item => (item.i === newLayoutItem.i ? newLayoutItem : item));
  const movingUp = y && layoutItem.y > y;

  // If this collides with anything, move it.
  // When doing this comparison, we have to sort the items we compare with
  // to ensure, in the case of multiple collisions, that we're getting the
  // nearest collision.
  const sorted = getSortedLayout(newLayout);

  if (movingUp) {
    sorted.reverse();
  }

  const collisions = getAllCollisions(sorted, newLayoutItem);

  // Move each item that collides away from this element.
  return collisions.reduce((acc, collision) => {
    if (collision.moved) {
      return acc;
    }

    // This makes it feel a bit more precise by waiting to swap for just a bit when moving up.
    if (newLayoutItem.y > collision.y && newLayoutItem.y - collision.y > collision.h / 4) {
      return acc;
    }

    if (isUserAction) {
      const fakeItem = {
        x: collision.x,
        w: collision.w,
        h: collision.h,
        i: '-1',
        y: Math.max(newLayoutItem.y - collision.h, 0),
      };

      if (!getFirstCollision(acc, fakeItem)) {
        return moveElement(acc, collision, collision.x, fakeItem.y, false);
      }
    }

    return moveElement(acc, collision, collision.x, collision.y + 1, false);
  }, newLayout);
}

/**
 * @param {GridLayout} compareWith
 * @param {GridLayoutItem} layoutItem
 * @returns {GridLayoutItem}
 */
const compactItem = (compareWith, layoutItem) => {
  const newLayoutItem = { ...layoutItem };

  while (newLayoutItem.y > 0 && !getFirstCollision(compareWith, newLayoutItem)) {
    newLayoutItem.y -= 1;
  }

  for (
    let collision = getFirstCollision(compareWith, newLayoutItem);
    collision;
    collision = getFirstCollision(compareWith, newLayoutItem)
  ) {
    newLayoutItem.y = collision.y + collision.h;
  }

  return newLayoutItem;
};

/**
 * @param {GridLayout} layout
 * @returns {GridLayout}
 */
const compact = (layout = []) => {
  const indexesById = layout.reduce((acc, { i }, index) => {
    acc[i] = index;

    return acc;
  }, {});

  const newLayout = new Array(layout.length);
  const sortedLayout = getSortedLayout(layout);
  const compareWith = [];

  for (const layoutItem of sortedLayout) {
    const newLayoutItem = compactItem(compareWith, layoutItem);

    newLayoutItem.moved = false;

    compareWith.push(newLayoutItem);

    newLayout[indexesById[layoutItem.i]] = newLayoutItem;
  }

  return newLayout;
};

export default {
  components: { CGridItem },
  mixins: [formArrayMixin],
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
    autoSize: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      placeholder: {
        x: 0,
        y: 0,
        w: 0,
        h: 0,
        i: -1,
      },
      width: null,
      resizing: false,
      moving: false,
      style: {},
    };
  },
  computed: {
    bind() {
      return {
        layout: this.layout,
        containerWidth: this.width,
        rowHeight: this.rowHeight,
        columnsCount: this.columnsCount,
        margin: this.margin,
        disabled: this.disabled,
      };
    },

    on() {
      return {
        resize: this.resizeItemHandler,
        resized: this.resizedItemHandler,
        move: this.moveItemHandler,
        moved: this.movedItemHandler,
      };
    },

    overviewLayoutStyle() {
      const rowHeightInPixels = `${WIDGET_GRID_ROW_HEIGHT}px`;

      return {
        padding: rowHeightInPixels,
        columnGap: rowHeightInPixels,
        gridTemplateColumns: `repeat(${WIDGET_GRID_COLUMNS_COUNT}, 1fr)`,
      };
    },

    preparedStyle() {
      return this.disabled ? this.overviewLayoutStyle : this.style;
    },
  },
  watch: {
    disabled(disabled) {
      if (disabled) {
        this.removeAllListeners();
        this.style = {};

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
    calculateContainerHeight() {
      if (!this.autoSize) {
        return '';
      }

      const [, marginY] = this.margin;
      const rowCount = calculateLayoutRowCount(this.layout);
      const bottom = calculateLayoutBottom(this.layout);

      return `${(bottom * this.rowHeight) + (marginY * rowCount) + marginY}px`;
    },

    updateHeight() {
      this.style = {
        height: this.calculateContainerHeight(),
      };
    },

    resizeObserverHandler() {
      const newWidth = this.$refs?.layout?.offsetWidth ?? null;

      if (newWidth !== this.width) {
        this.width = newWidth;
      }

      this.updateHeight();
    },

    resizeItemHandler(id, x, y, h, w) {
      const index = this.layout.findIndex(item => item.i === id);
      const newLayoutItem = { ...this.layout[index], h, w };
      const newLayout = compact(this.layout.map(item => (item.i === id ? newLayoutItem : item)));

      this.resizing = true;

      this.placeholder = {
        ...this.placeholder,
        ...newLayout.find(item => item.i === id),
      };

      this.updateModel(newLayout);

      this.$nextTick(() => {
        this.updateHeight();
      });
    },

    resizedItemHandler(id, x, y, h, w) {
      const index = this.layout.findIndex(item => item.i === id);
      const newLayoutItem = { ...this.layout[index], h, w };

      this.updateModel(compact(this.layout.map(item => (item.i === id ? newLayoutItem : item))));

      this.$nextTick(() => {
        this.updateHeight();
        this.resizing = false;
      });
    },

    moveItemHandler(id, x, y) {
      this.moving = true;
      const layoutItem = this.layout.find(item => item.i === id);
      const newLayout = compact(moveElement(this.layout, layoutItem, x, y, true));
      this.placeholder = {
        ...this.placeholder,
        ...newLayout.find(item => item.i === id),
      };

      this.updateModel(newLayout);

      this.$nextTick(() => {
        this.updateHeight();
      });
    },

    movedItemHandler(id, x, y) {
      const layoutItem = this.layout.find(item => item.i === id);

      this.updateModel(compact(moveElement(this.layout, layoutItem, x, y, true)));

      this.$nextTick(() => {
        this.moving = false;
        this.updateHeight();
      });
    },

    addAllListeners() {
      this.$resizeObserver = new ResizeObserver(this.resizeObserverHandler);
      this.$resizeObserver.observe(this.$el);
    },

    removeAllListeners() {
      this.$resizeObserver.disconnect();
    },
  },
};
</script>

<style lang="scss" scoped>
.c-grid-layout {
  &:not(&--disabled) {
    background-color: rgba(60, 60, 60, 0.05);
    margin: auto;
    position: relative;
  }

  &__placeholder {
    background: var(--v-primary-darken1);
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
