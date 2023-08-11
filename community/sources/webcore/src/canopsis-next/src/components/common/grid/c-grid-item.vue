<template lang="pug">
  div.c-grid-item(:style="style")
    slot
    div.c-grid-item__resize-handler(ref="resizeHandler")
      v-icon(small) $vuetify.icons.resize_right
</template>

<script>
import { debounce } from 'lodash';

/**
 * Get count above grid-item
 *
 * @param {number} cardX
 * @param {number} cardY
 * @param {number} currentCardWidth
 * @returns {number}
 */
const getCountAboveItems = (layout = [], itemX, itemY, itemW) => {
  const { count } = layout
    .filter(({ y }) => y < itemY)
    .sort((a, b) => b.y - a.y)
    .reduce((acc, {
      y, x, w, h,
    }) => {
      if (acc.y !== y + h) {
        return acc;
      }

      const range = (acc.x + acc.w) - (x + w);

      const isInteraction = range > 0 ? range < acc.w : Math.abs(range) < w;

      if (y < acc.y && isInteraction) {
        acc.count += 1;
        acc.x = x;
        acc.y = y;
        acc.w = w;
      }

      return acc;
    }, {
      count: 0,
      x: itemX,
      y: itemY,
      w: itemW,
    });

  return count;
};

export default {
  props: {
    layout: {
      type: Array,
      required: true,
    },
    i: {
      type: [Number, String],
      required: true,
    },
    containerWidth: {
      type: Number,
      default: 100,
    },
    rowHeight: {
      type: Number,
      default: 150,
    },
    columnsCount: {
      type: Number,
      default: 1,
    },
    margin: {
      type: Array,
      default: () => [10, 10],
    },
    x: {
      type: Number,
      default: 0,
    },
    y: {
      type: Number,
      default: 0,
    },
    w: {
      type: Number,
      default: 0,
    },
    h: {
      type: Number,
      default: 0,
    },
    minW: {
      type: Number,
      default: 0,
    },
    maxW: {
      type: Number,
      default: Infinity,
    },
    minH: {
      type: Number,
      default: 0,
    },
    maxH: {
      type: Number,
      default: Infinity,
    },
    maxRows: {
      type: Number,
      default: Infinity,
    },
    autoHeight: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    overlay: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      resizable: null,
      draggable: null,
      isResizing: false,
      isDragging: false,
      resizing: null,
      dragging: null,
      lastX: NaN,
      lastY: NaN,
      lastW: NaN,
      lastH: NaN,
      style: {},

      previousW: null,
      previousH: null,
      previousX: null,
      previousY: null,
    };
  },
  computed: {
    columnWidth() {
      return (this.containerWidth - (this.margin[0] * (this.columnsCount + 1))) / this.columnsCount;
    },

    innerX() {
      return this.x + this.w > this.columnsCount
        ? 0
        : this.x;
    },

    innerW() {
      return (this.x + this.w > this.columnsCount) && (this.w > this.columnsCount)
        ? this.columnsCount
        : this.w;
    },
  },
  mounted() {
    this.calculateStyle();
    this.addAllAutoSizeListeners();
    this.addAllMainPropsWatchers();
  },
  beforeDestroy() {
    this.removeAllAutoSizeListeners();
    this.removeAllMainPropsWatchers();
  },
  methods: {
    addAllMainPropsWatchers() {
      this.$mainPropsWatchers = [
        'x',
        'y',
        'w',
        'h',
        'rowHeight',
        'columnsCount',
      ].map(prop => this.$watch(prop, this.calculateStyle));

      this.$mainPropsWatchers.push(
        this.$watch('containerWidth', () => {
          this.calculateStyle();

          if (this.autoHeight) {
            this.callAutoSizeHeight();
          }
        }),

        this.$watch('autoHeight', () => {
          if (this.autoHeight) {
            this.addAllAutoSizeListeners();
          } else {
            this.removeAllAutoSizeListeners();
          }
        }),
      );
    },

    removeAllMainPropsWatchers() {
      if (!this.$mainPropsWatchers?.length) {
        return;
      }

      this.$mainPropsWatchers.forEach(unwatch => unwatch());

      delete this.$mainPropsWatchers;
    },

    addAllAutoSizeListeners() {
      if (!this.autoHeight) {
        return;
      }

      /**
       * Method which wrap autoSizeHeight method into setTimeout and debounce
       * We are using debounce here for pauses between calls
       * And setTimeout instead of $nextTick for correct call after rerender
       * We are declared that not in the methods because we need one method for every instance, not for all
       */
      this.callAutoSizeHeight = debounce(() => {
        setTimeout(() => {
          this.autoSizeHeight();
        }, 0);
      }, 100);

      this.callAutoSizeHeight();

      this.$mutationObserver = new MutationObserver(this.callAutoSizeHeight);

      /**
       * We are observe on the mutationObserver after render
       */
      this.$nextTick(() => {
        const element = this.getDefaultSlotElement();

        if (!element) {
          return;
        }

        this.$mutationObserver.observe(element, {
          childList: true,
          subtree: true,
        });
      });
    },

    removeAllAutoSizeListeners() {
      if (!this.$mutationObserver) {
        return;
      }

      this.$mutationObserver.disconnect();

      delete this.$mutationObserver;
      delete this.callAutoSizeHeight;
    },

    /**
     * Get defaultSlot element
     *
     * @returns {null|HTMLElement}
     */
    getDefaultSlotElement() {
      const { default: defaultSlots } = this.$slots;

      if (!defaultSlots) {
        return null;
      }

      const [{ elm: element }] = defaultSlots;

      return element || null;
    },

    calculateWH(height, width) {
      const [marginX] = this.margin;
      const w = Math.round((width + marginX) / (this.columnWidth + marginX));
      const h = Math.round((height) / this.rowHeight);

      return {
        w: Math.min(Math.max(Math.min(w, this.columnsCount - this.innerX), this.minW), this.maxW),
        h: Math.min(Math.max(Math.min(h, this.maxRows - this.y), this.minH), this.maxH),
      };
    },

    calculatePosition(x, y, w, h) {
      const [marginX, marginY] = this.margin;
      const marginCount = getCountAboveItems(this.layout, this.x, this.y, this.w);

      return {
        left: Math.round(this.columnWidth * x + (x + 1) * marginX),
        top: Math.round((this.rowHeight * y) + (marginY * marginCount) + marginY),
        width: w === Infinity ? w : Math.round(this.columnWidth * w + Math.max(0, w - 1) * marginX),
        height: h === Infinity ? h : Math.round(this.rowHeight * h),
      };
    },

    calculateStyle() {
      const pos = this.calculatePosition(this.innerX, this.y, this.innerW, this.h);

      if (this.isDragging) {
        pos.top = this.dragging.top;
        pos.left = this.dragging.left;
      }

      if (this.isResizing) {
        pos.width = this.resizing.width;
        pos.height = this.resizing.height;
      }

      const translate = `translate3d(${pos.left}px,${pos.top}px, 0)`;

      this.style = {
        transform: translate,
        WebkitTransform: translate,
        MozTransform: translate,
        msTransform: translate,
        OTransform: translate,
        width: `${pos.width}px`,
        height: `${pos.height}px`,
        position: 'absolute',
      };
    },

    autoSizeHeight() {
      this.previousW = this.innerW;
      this.previousH = this.h;
      const element = this.getDefaultSlotElement();

      if (!element) {
        return;
      }

      const newSize = element.getBoundingClientRect();
      const pos = this.calculateWH(newSize.height, newSize.width);

      pos.h = Math.max(Math.min(pos.h, this.maxH), this.minH);

      if (this.h !== pos.h) {
        this.$emit('resize', this.i, pos.h, this.innerW, newSize.height, newSize.width);
      }

      if (this.previousH !== pos.h) {
        this.$emit('resized', this.i, this.innerX, this.y, pos.h, this.innerW);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.c-grid-item {
  position: absolute;
  height: 1000px;
  overflow: hidden;
  transition: none;
  left: 0;
  right: auto;

  &__resize-handler {
    cursor: se-resize;
    position: absolute;
    width: 16px;
    height: 16px;
    bottom: 3px;
    right: 3px;
    z-index: 3;
    opacity: .7;
  }

  &__placeholder {
    background: red;
    opacity: 0.2;
    transition-duration: 100ms;
    z-index: 2;
    user-select: none;
  }

  &:after {
    content: '';
    background-color: #888;
    position: absolute;
    left: 0;
    top: 36px;
    width: 100%;
    height: 100%;
    opacity: 0.4;
    z-index: 2;
  }
}
</style>
