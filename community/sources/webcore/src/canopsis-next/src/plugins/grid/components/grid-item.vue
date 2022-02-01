<script>
import { debounce } from 'lodash';
import { GridItem } from 'vue-grid-layout';
import { createCoreData, getControlPosition } from '../helpers/draggable-utils';

export default {
  extends: GridItem,
  props: {
    autoHeight: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    autoHeight(value) {
      if (value) {
        this.subscribeToAllChangesForAutoSizeHeight();
      } else {
        this.unsubscribeFromAllChangesForAutoSizeHeight();
      }
    },
  },
  beforeCreate() {
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

    this.$_mutationObserver = new MutationObserver(this.callAutoSizeHeight);
  },
  mounted() {
    if (this.autoHeight) {
      this.subscribeToAllChangesForAutoSizeHeight();
    }
  },
  beforeDestroy() {
    this.unsubscribeFromAllChangesForAutoSizeHeight();
  },
  methods: {
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

    /**
     * Call callAutoSizeHeight method and subscribe to 'windowResizeEvent' and observe defaultSlot element
     * For callAutoSizeHeight method calling
     */
    subscribeToAllChangesForAutoSizeHeight() {
      this.callAutoSizeHeight();
      this.eventBus.$on('windowResizeEvent', this.callAutoSizeHeight);

      /**
       * We are observe on the mutationObserver after render
       */
      this.$nextTick(() => {
        const element = this.getDefaultSlotElement();

        if (!element) {
          return;
        }

        this.$_mutationObserver.observe(element, {
          childList: true,
          subtree: true,
        });
      });
    },

    /**
     * Unsubscribe from 'windowResizeEvent' and disconnect from defaultSlot element observing
     */
    unsubscribeFromAllChangesForAutoSizeHeight() {
      this.eventBus.$off('windowResizeEvent', this.callAutoSizeHeight);
      this.$_mutationObserver.disconnect();
    },

    /**
     * Calculate grid-item height size by defaultSlot element height size
     */
    autoSizeHeight() {
      this.previousW = this.innerW;
      this.previousH = this.innerH;
      const element = this.getDefaultSlotElement();

      if (!element) {
        return;
      }

      const newSize = element.getBoundingClientRect();
      const pos = this.calcWH(newSize.height, newSize.width);

      if (pos.h < this.minH) {
        pos.h = this.minH;
      }
      if (pos.h > this.maxH) {
        pos.h = this.maxH;
      }

      if (pos.h < 1) {
        pos.h = 1;
      }

      if (this.innerH !== pos.h) {
        this.$emit('resize', this.i, pos.h, this.innerW, newSize.height, newSize.width);
      }

      if (this.previousH !== pos.h) {
        this.$emit('resized', this.i, pos.h, this.previousW, newSize.height, newSize.width);
        this.eventBus.$emit('resizeEvent', 'resizeend', this.i, this.innerX, this.innerY, pos.h, this.innerW);
      }
    },

    /**
     * Handle resize grid-item method
     *
     * @param {Event} event
     */
    handleResize(event) {
      if (this.static) {
        return;
      }

      const element = this.getDefaultSlotElement();
      const position = getControlPosition(event);

      if (!position || !element) {
        return;
      }

      const { x, y } = position;
      const newSize = { width: 0, height: 0 };

      let pos;

      switch (event.type) {
        case 'resizestart':
          this.previousW = this.innerW;
          this.previousH = this.innerH;

          pos = this.calcPosition(this.innerX, this.innerY, this.innerW, this.innerH);

          newSize.width = pos.width;
          newSize.height = pos.height;

          this.resizing = newSize;
          this.isResizing = true;

          break;
        case 'resizemove': {
          const newElementSize = element.getBoundingClientRect();
          const coreEvent = createCoreData(this.lastW, this.lastH, x, y);

          if (this.renderRtl) {
            newSize.width = this.resizing.width - coreEvent.deltaX;
          } else {
            newSize.width = this.resizing.width + coreEvent.deltaX;
          }

          newSize.height = this.autoHeight
            ? newElementSize.height
            : this.resizing.height + coreEvent.deltaY;

          this.resizing = newSize;
          break;
        }
        case 'resizeend': {
          pos = this.calcPosition(this.innerX, this.innerY, this.innerW, this.innerH);

          newSize.width = pos.width;
          newSize.height = pos.height;

          this.resizing = null;
          this.isResizing = false;
          break;
        }
      }

      // Get new WH
      pos = this.calcWH(newSize.height, newSize.width);
      if (pos.w < this.minW) {
        pos.w = this.minW;
      }
      if (pos.w > this.maxW) {
        pos.w = this.maxW;
      }
      if (pos.h < this.minH) {
        pos.h = this.minH;
      }
      if (pos.h > this.maxH) {
        pos.h = this.maxH;
      }

      if (pos.h < 1) {
        pos.h = 1;
      }
      if (pos.w < 1) {
        pos.w = 1;
      }

      this.lastW = x;
      this.lastH = y;

      if (this.innerW !== pos.w || this.innerH !== pos.h) {
        this.$emit('resize', this.i, pos.h, pos.w, newSize.height, newSize.width);
      }
      if (event.type === 'resizeend' && (this.previousW !== this.innerW || this.previousH !== this.innerH)) {
        this.$emit('resized', this.i, pos.h, pos.w, newSize.height, newSize.width);

        if (this.autoHeight) {
          this.callAutoSizeHeight();
        }
      }

      this.eventBus.$emit('resizeEvent', event.type, this.i, this.innerX, this.innerY, pos.h, pos.w);
    },

    /**
     * Get count above grid-item
     *
     * @param {number} currentCardX
     * @param {number} currentCardY
     * @param {number} currentCardWidth
     * @returns {number}
     */
    countAboveCard(currentCardX, currentCardY, currentCardWidth) {
      const { layout } = this.$parent;

      const { count } = layout
        .filter(({ y }) => y < currentCardY)
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
          x: currentCardX,
          y: currentCardY,
          w: currentCardWidth,
        });

      return count;
    },

    /**
     * Calculate grid-item position
     *
     * @param {number} x
     * @param {number} y
     * @param {number} w
     * @param {number} h
     * @returns {Object}
     */
    calcPosition(x, y, w, h) {
      const [marginX, marginY] = this.margin;
      const colWidth = this.calcColWidth();
      const height = Math.round(this.rowHeight * h);

      const marginCount = this.countAboveCard(x, y, w);

      const result = {
        top: Math.round((this.rowHeight * y) + (marginY * marginCount) + marginY),
        width: w === Infinity ? w : Math.round((colWidth * w) + (Math.max(0, w - 1) * marginX)),
        height: h === Infinity ? h : height,
      };

      const horizontal = Math.round((colWidth * x) + ((x + 1) * marginX));

      if (this.renderRtl) {
        result.right = horizontal;
      } else {
        result.left = horizontal;
      }

      return result;
    },

    /**
     * Calculate width and height in grid parameters
     *
     * @param {number} height
     * @param {number} width
     * @returns {{w: number, h: number}}
     */
    calcWH(height, width) {
      const [marginX] = this.margin;
      const colWidth = this.calcColWidth();

      let w = Math.round((width + marginX) / (colWidth + marginX));
      let h = Math.round((height) / this.rowHeight);

      w = Math.max(Math.min(w, this.cols - this.innerX), 0);
      h = Math.max(Math.min(h, this.maxRows - this.innerY), 0);

      return { w, h };
    },

    /**
     * Calculate x and y positions in grid parameters
     *
     * @param {number} top
     * @param {number} left
     * @returns {{x: number, y: number}}
     */
    calcXY(top, left) {
      const [marginX] = this.margin;
      const colWidth = this.calcColWidth();

      let x = Math.round((left - marginX) / (colWidth + marginX));
      let y = Math.round(top / this.rowHeight);

      x = Math.max(Math.min(x, this.cols - this.innerW), 0);
      y = Math.max(Math.min(y, this.maxRows - this.innerH), 0);

      return { x, y };
    },
  },
};
</script>
