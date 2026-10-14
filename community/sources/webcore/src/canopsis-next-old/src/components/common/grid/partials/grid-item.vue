<template>
  <v-card
    :class="cardClass"
    :style="style"
    class="grid-item"
  >
    <div ref="defaultSlotWrapper">
      <slot
        :on="slotOn"
        :item="item"
      />
    </div>
    <div
      v-if="!disabled"
      class="grid-item__resize-handler"
      v-on="resizeHandlerOn"
    >
      <v-icon small>
        $vuetify.icons.resize_right
      </v-icon>
    </div>
  </v-card>
</template>

<script>
import { debounce, throttle } from 'lodash';

import { getControlPosition, getCountAboveItems, getItemDelta } from '@/helpers/grid';

export default {
  props: {
    layout: {
      type: Array,
      required: true,
    },
    item: {
      type: Object,
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
    disabled: {
      type: Boolean,
      default: false,
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
      resizing: null,
      dragging: null,
    };
  },
  computed: {
    cardClass() {
      return {
        'grid-item--resizing': this.resizing || this.dragging,
        'grid-item--disabled': this.disabled,
      };
    },

    columnWidth() {
      return (this.containerWidth - (this.margin[0] * (this.columnsCount + 1))) / this.columnsCount;
    },

    slotOn() {
      if (this.disabled) {
        return {};
      }

      return {
        pointerdown: this.dragStartHandler,
        toggle: this.toggleAutoHeightHandler,
      };
    },

    resizeHandlerOn() {
      return {
        pointerdown: this.resizeStartHandler,
      };
    },

    x() {
      const { x = 0, w = 1 } = this.item;

      return x + w > this.columnsCount
        ? 0
        : x;
    },

    y() {
      return this.item.y ?? 0;
    },

    w() {
      const { x = 0, w = 1 } = this.item;

      return (x + w > this.columnsCount) && (w > this.columnsCount)
        ? this.columnsCount
        : w;
    },

    h() {
      return this.item.h ?? 1;
    },

    autoHeight() {
      return this.item.autoHeight ?? false;
    },

    layoutElement() {
      return this.$parent?.$el ?? document.body;
    },

    disabledItemStyle() {
      const [, marginY = 0] = this.margin;
      return {
        margin: `${marginY / 2}px 0`,
        gridColumnStart: this.x + 1,
        gridColumnEnd: this.x + 1 + this.w,
        gridRowStart: this.y + 1,
        gridRowEnd: this.y + this.h + 1,
        height: this.autoHeight ? 'auto' : `${this.rowHeight * this.h}px`,
      };
    },

    enabledItemStyle() {
      const position = this.calculatePosition(this.x, this.y, this.w, this.h);

      if (this.dragging) {
        position.top = this.dragging.top;
        position.left = this.dragging.left;
      }

      if (this.resizing) {
        position.width = this.resizing.width;
        position.height = this.resizing.height;
      }

      const translate = `translate3d(${position.left}px,${position.top}px, 0)`;

      return {
        transform: translate,
        WebkitTransform: translate,
        MozTransform: translate,
        msTransform: translate,
        OTransform: translate,
        width: `${position.width}px`,
        height: `${position.height}px`,
        position: 'absolute',
      };
    },

    style() {
      return this.disabled ? this.disabledItemStyle : this.enabledItemStyle;
    },
  },
  watch: {
    disabled(disabled) {
      if (disabled) {
        this.removeAllWatchersAndListeners();
        return;
      }

      this.addAllWatchersAndListeners();
    },
  },
  created() {
    this.throttledResizeProcessHandler = throttle(this.resizeProcessHandler, this.throttle);
    this.throttledDragProcessHandler = throttle(this.dragProcessHandler, this.throttle);
    this.debouncedSetAutoCalculatedHeight = debounce(() => setTimeout(this.setAutoCalculatedHeight, 0), this.debounce);
  },
  mounted() {
    if (!this.disabled) {
      this.addAllWatchersAndListeners();
    }
  },
  beforeDestroy() {
    this.removeAllWatchersAndListeners();
  },
  methods: {
    toggleAutoHeightHandler() {
      this.$emit('toggle', this.item);
    },

    addAllWatchersAndListeners() {
      this.addAllAutoHeightListeners();
      this.addAllMainPropsWatchers();
    },

    removeAllWatchersAndListeners() {
      this.removeAllAutoHeightListeners();
      this.removeAllMainPropsWatchers();
      this.removeAllResizeListeners();
      this.removeAllDragListeners();
    },

    addAllMainPropsWatchers() {
      if (this.$mainPropsWatchers) {
        return;
      }

      this.$mainPropsWatchers = [
        this.$watch('containerWidth', () => {
          if (this.autoHeight) {
            this.debouncedSetAutoCalculatedHeight();
          }
        }),

        this.$watch('autoHeight', () => {
          if (this.autoHeight) {
            this.addAllAutoHeightListeners();
          } else {
            this.removeAllAutoHeightListeners();
          }
        }),
      ];
    },

    removeAllMainPropsWatchers() {
      if (!this.$mainPropsWatchers?.length) {
        return;
      }

      this.$mainPropsWatchers.forEach(unwatch => unwatch());

      delete this.$mainPropsWatchers;
    },

    calculateXY(top, left) {
      const [marginX] = this.margin;

      const x = Math.round((left - marginX) / (this.columnWidth + marginX));
      const y = Math.round(top / this.rowHeight);

      return {
        x: Math.max(Math.min(x, this.columnsCount - this.w), 0),
        y: Math.max(y, 0),
      };
    },

    calculateWH(width, height) {
      const [marginX] = this.margin;
      const w = Math.round((width + marginX) / (this.columnWidth + marginX));
      const h = Math.round((height) / this.rowHeight);

      return {
        w: Math.min(Math.max(Math.min(w, this.columnsCount - this.x), this.minW), this.maxW),
        h: Math.min(Math.max(h, this.minH), this.maxH),
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

    /**
     * AUTO HEIGHT
     */
    addAllAutoHeightListeners() {
      if (!this.autoHeight) {
        return;
      }

      this.debouncedSetAutoCalculatedHeight();

      this.$mutationObserver = new MutationObserver(this.debouncedSetAutoCalculatedHeight);

      /**
       * We are observe on the mutationObserver after render
       */
      this.$nextTick(() => {
        const element = this.$refs.defaultSlotWrapper;

        if (!element) {
          return;
        }

        this.$mutationObserver.observe(element, {
          childList: true,
          subtree: true,
        });
      });
    },

    removeAllAutoHeightListeners() {
      if (!this.$mutationObserver) {
        return;
      }

      this.$mutationObserver.disconnect();

      delete this.$mutationObserver;
    },

    setAutoCalculatedHeight() {
      const element = this.$refs.defaultSlotWrapper;

      if (!element) {
        return;
      }

      const { height = 0, width = 0 } = element.getBoundingClientRect();
      let { h } = this.calculateWH(width, height);

      h = Math.max(Math.min(h, this.maxH), this.minH);

      if (this.h !== h) {
        this.$emit('resized', this.item, this.w, h);
      }
    },

    /**
     * RESIZING
     */
    removeAllResizeListeners() {
      window.removeEventListener('pointerup', this.resizeEndHandler);
      window.removeEventListener('pointermove', this.throttledResizeProcessHandler);
    },

    setWHByResizing(position, controlPosition) {
      this.resizing = position;
      this.controlX = controlPosition.x;
      this.controlY = controlPosition.y;

      const { w, h } = this.calculateWH(position.width, position.height);

      if (this.w !== w || this.h !== h) {
        this.$emit('resize', this.item, w, h);
      }
    },

    resizeStartHandler(event) {
      if (event.button) {
        return;
      }

      window.addEventListener('pointerup', this.resizeEndHandler);
      window.addEventListener('pointermove', this.throttledResizeProcessHandler);

      this.setWHByResizing(
        this.calculatePosition(this.x, this.y, this.w, this.h),
        getControlPosition(event, this.layoutElement),
      );
    },

    resizeProcessHandler(event) {
      const element = this.$refs.defaultSlotWrapper;

      if (!this.resizing || !element) {
        return;
      }

      const { height = 0 } = element.getBoundingClientRect();
      const controlPosition = getControlPosition(event, this.layoutElement);
      const delta = getItemDelta(this.controlX, this.controlY, controlPosition.x, controlPosition.y);

      this.setWHByResizing(
        {
          width: this.resizing.width + delta.deltaX,
          height: this.autoHeight
            ? height
            : this.resizing.height + delta.deltaY,
        },
        controlPosition,
      );
    },

    resizeEndHandler() {
      if (!this.resizing) {
        return;
      }

      window.removeEventListener('pointermove', this.throttledResizeProcessHandler);

      this.resizing = this.calculatePosition(this.x, this.y, this.w, this.h);

      const { w, h } = this.calculateWH(this.resizing.width, this.resizing.height);

      this.$emit('resized', this.item, w, h);

      if (this.autoHeight) {
        this.debouncedSetAutoCalculatedHeight();
      }

      this.resizing = false;
    },

    /**
     * DRAG AND DROP
     */
    removeAllDragListeners() {
      window.removeEventListener('pointermove', this.throttledDragProcessHandler);
      window.removeEventListener('pointerup', this.dragEndHandler);
    },

    setXYByDragging(position, controlPosition) {
      this.dragging = position;
      this.controlX = controlPosition.x;
      this.controlY = controlPosition.y;

      const { x, y } = this.calculateXY(position.top, position.left);

      if (this.x !== x || this.y !== y) {
        this.$emit('move', this.item, x, y);
      }
    },

    dragStartHandler(event) {
      if (this.resizing) {
        return;
      }

      window.addEventListener('pointermove', this.throttledDragProcessHandler);
      window.addEventListener('pointerup', this.dragEndHandler);

      const parentRect = this.layoutElement.getBoundingClientRect();
      const clientRect = event.target.getBoundingClientRect();

      this.setXYByDragging(
        {
          top: clientRect.top - parentRect.top,
          left: clientRect.left - parentRect.left,
        },
        getControlPosition(event, this.layoutElement),
      );
    },

    dragProcessHandler(event) {
      if (!this.dragging || this.resizing) {
        return;
      }

      const controlPosition = getControlPosition(event, this.layoutElement);
      const delta = getItemDelta(this.controlX, this.controlY, controlPosition.x, controlPosition.y);

      this.setXYByDragging(
        {
          top: this.dragging.top + delta.deltaY,
          left: this.dragging.left + delta.deltaX,
        },
        controlPosition,
      );
    },

    dragEndHandler(event) {
      if (!this.dragging || this.resizing) {
        return;
      }

      window.removeEventListener('pointermove', this.throttledDragProcessHandler);
      window.removeEventListener('pointerup', this.dragEndHandler);

      const controlPosition = getControlPosition(event, this.layoutElement);
      const delta = getItemDelta(this.controlX, this.controlY, controlPosition.x, controlPosition.y);
      const { x, y } = this.calculateXY(this.dragging.top + delta.deltaY, this.dragging.left + delta.deltaX);

      this.dragging = null;

      this.$emit('moved', this.item, x, y);
    },
  },
};
</script>

<style lang="scss" scoped>
.grid-item {
  transition: none;
  overflow: auto;

  &__resize-handler {
    display: none;
    cursor: se-resize;
    position: absolute;
    width: 16px;
    height: 16px;
    line-height: 16px;
    bottom: 3px;
    right: 3px;
    z-index: 3;
    opacity: .7;
  }

  &:not(&--disabled) {
    position: absolute;
    left: 0;
    right: auto;
    z-index: 2;
    overflow: hidden;
    height: 1000px;

    &.grid-item--resizing {
      opacity: .7;
      z-index: 999;
      user-select: none;
    }

    .grid-item__resize-handler {
      display: block;
    }

    &:after {
      content: '';
      background-color: #888;
      position: absolute;
      left: 0;
      top: 0;
      width: 100%;
      height: 100%;
      opacity: 0.3;
      z-index: 2;
    }
  }
}
</style>
