<template lang="pug">
  v-card.c-grid-item(
    :class="{ 'c-grid-item--resizing': resizing || dragging, 'c-grid-item--disabled': disabled }",
    :style="preparedStyle",
    :draggable="!disabled"
  )
    slot
    div.c-grid-item__resize-handler(v-if="resizable", ref="resizeHandler")
      v-icon(small) $vuetify.icons.resize_right
</template>

<script>
import { debounce, isNumber, throttle } from 'lodash';

import { MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP, WIDGET_GRID_ROW_HEIGHT } from '@/constants';

/**
 * TODO: MOVE TO HELPERS
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

export const createCoreData = (lastX, lastY, x, y) => (
  !isNumber(lastX)
    ? {
      deltaX: 0,
      deltaY: 0,
      lastX: x,
      lastY: y,
      x,
      y,
    }
    : {
      deltaX: x - lastX,
      deltaY: y - lastY,
      lastX,
      lastY,
      x,
      y,
    }
);

export const getControlPosition = (event, layoutElement) => {
  const layoutElementRect = layoutElement.getBoundingClientRect();

  const x = event.clientX + layoutElement.scrollLeft - layoutElementRect.left;
  const y = event.clientY + layoutElement.scrollTop - layoutElementRect.top;

  return { x, y };
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
    widget: {
      type: Object,
      default: () => ({}),
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
    resizable: {
      type: Boolean,
      default: false,
    },
    throttle: {
      type: Number,
      default: 50,
    },
    debounce: {
      type: Number,
      default: 100,
    },
  },
  data() {
    return {
      resizing: null,
      dragging: null,
      lastX: NaN,
      lastY: NaN,
      lastW: NaN,
      lastH: NaN,
      style: {},
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

    layoutElement() {
      return this.$parent?.$el ?? document.body;
    },

    gridParameters() {
      const key = MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP[this.$mq];

      return this.widget.grid_parameters[key];
    },

    overviewItemStyle() {
      return {
        margin: `${WIDGET_GRID_ROW_HEIGHT / 2}px 0`,
        gridColumnStart: this.gridParameters.x + 1,
        gridColumnEnd: this.gridParameters.x + 1 + this.gridParameters.w,
        gridRowStart: this.gridParameters.y + 1,
        gridRowEnd: this.gridParameters.y + this.gridParameters.h + 1,
        height: this.gridParameters.autoHeight ? 'auto' : `${WIDGET_GRID_ROW_HEIGHT * this.gridParameters.h}px`,
      };
    },

    preparedStyle() {
      return this.disabled ? this.overviewItemStyle : this.style;
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
    this.throttledMousemoveHandler = throttle(this.mousemoveHandler, this.throttle);
    this.throttledDragoverHandler = throttle(this.dragoverHandler, this.throttle);
    this.debouncedAutoSizeHeight = debounce(() => setTimeout(this.autoSizeHeight, 0), this.debounce);
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
    addAllWatchersAndListeners() {
      this.calculateStyle();
      this.addAllAutoSizeListeners();
      this.addAllMainPropsWatchers();
      this.addAllResizeListeners();
      this.addAllDragListeners();
    },

    removeAllWatchersAndListeners() {
      this.removeAllAutoSizeListeners();
      this.removeAllMainPropsWatchers();
      this.removeAllResizeListeners();
      this.removeAllDragListeners();
    },

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
            this.debouncedAutoSizeHeight();
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

    calculateWH(height, width) {
      const [marginX] = this.margin;
      const w = Math.round((width + marginX) / (this.columnWidth + marginX));
      const h = Math.round((height) / this.rowHeight);

      return {
        w: Math.min(Math.max(Math.min(w, this.columnsCount - this.innerX), this.minW), this.maxW),
        h: Math.min(Math.max(Math.min(h, this.maxRows - this.y), this.minH), this.maxH),
      };
    },

    calculateXY(top, left) {
      const [marginX] = this.margin;

      let x = Math.round((left - marginX) / (this.columnWidth + marginX));
      let y = Math.round(top / this.rowHeight);

      x = Math.max(Math.min(x, this.columnsCount - this.innerW), 0);
      y = Math.max(Math.min(y, this.maxRows - this.h), 0);

      return { x, y };
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

      if (this.dragging) {
        pos.top = this.dragging.top;
        pos.left = this.dragging.left;
      }

      if (this.resizing) {
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

    /**
     * AUTO SIZE
     */
    addAllAutoSizeListeners() {
      if (!this.autoHeight) {
        return;
      }

      this.debouncedAutoSizeHeight();

      this.$mutationObserver = new MutationObserver(this.debouncedAutoSizeHeight);

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
    },

    getDefaultSlotElement() {
      const { default: defaultSlots } = this.$slots;

      if (!defaultSlots) {
        return null;
      }

      const [{ elm: element }] = defaultSlots;

      return element || null;
    },

    autoSizeHeight() {
      const element = this.getDefaultSlotElement();

      if (!element) {
        return;
      }

      const newSize = element.getBoundingClientRect();
      const pos = this.calculateWH(newSize.height, newSize.width);

      pos.h = Math.max(Math.min(pos.h, this.maxH), this.minH);

      if (this.h !== pos.h) {
        this.$emit('resized', this.i, this.innerX, this.y, pos.h, this.innerW);
      }
    },

    /**
     * RESIZING
     */
    addAllResizeListeners() {
      this.$refs.resizeHandler?.addEventListener('mousedown', this.mousedownHandler);
      document.addEventListener('mouseup', this.mouseupHandler);
    },

    removeAllResizeListeners() {
      this.$refs.resizeHandler?.removeEventListener('mousedown', this.mousedownHandler);
      document.removeEventListener('mouseup', this.mouseupHandler);
    },

    mousedownHandler(event) {
      document.addEventListener('mousemove', this.throttledMousemoveHandler);

      const position = getControlPosition(event, this.layoutElement);
      this.resizing = this.calculatePosition(this.innerX, this.y, this.innerW, this.h);

      const { w, h } = this.calculateWH(this.resizing.height, this.resizing.width);

      this.$emit('resize', this.i, this.x, this.y, h, w);

      this.lastW = position.x;
      this.lastH = position.y;

      this.calculateStyle();
    },

    mousemoveHandler(event) {
      if (!this.resizing) {
        return;
      }

      const position = getControlPosition(event, this.layoutElement);
      const element = this.getDefaultSlotElement();
      const newElementSize = element.getBoundingClientRect();
      const coreEvent = createCoreData(this.lastW, this.lastH, position.x, position.y);

      this.lastW = position.x;
      this.lastH = position.y;
      this.resizing = {
        width: this.resizing.width + coreEvent.deltaX,
        height: this.autoHeight
          ? newElementSize.height
          : this.resizing.height + coreEvent.deltaY,
      };

      const { w, h } = this.calculateWH(this.resizing.height, this.resizing.width);

      if (this.innerW !== w || this.h !== h) {
        this.$emit('resize', this.i, this.x, this.y, h, w);
      }

      this.calculateStyle();
    },

    mouseupHandler() {
      if (!this.resizing) {
        return;
      }

      document.removeEventListener('mousemove', this.throttledMousemoveHandler);

      this.resizing = this.calculatePosition(this.innerX, this.y, this.innerW, this.h);

      const { w, h } = this.calculateWH(this.resizing.height, this.resizing.width);

      this.$emit('resized', this.i, this.x, this.y, h, w);

      if (this.autoHeight) {
        this.debouncedAutoSizeHeight();
      }

      this.resizing = false;
      this.calculateStyle();
    },

    /**
     * DRAG AND DROP
     */
    addAllDragListeners() {
      this.$el.addEventListener('dragstart', this.dragstartHandler);
    },

    removeAllDragListeners() {
      this.$el.removeEventListener('dragstart', this.dragstartHandler);
      this.$el.removeEventListener('drag', this.throttledDragoverHandler);
      this.$el.removeEventListener('dragend', this.dragendHandler);
    },

    dragstartHandler(event) {
      if (this.resizing) {
        event.preventDefault();

        return;
      }

      this.layoutElement.addEventListener('dragover', this.throttledDragoverHandler);
      this.$el.addEventListener('dragend', this.dragendHandler);

      const position = getControlPosition(event, this.layoutElement);
      const newPosition = { top: 0, left: 0 };
      const { x, y } = position;

      const parentRect = this.layoutElement.getBoundingClientRect();
      const clientRect = event.target.getBoundingClientRect();
      newPosition.left = clientRect.left - parentRect.left;
      newPosition.top = clientRect.top - parentRect.top;
      this.dragging = newPosition;
      const pos = this.calculateXY(newPosition.top, newPosition.left);

      this.lastX = x;
      this.lastY = y;

      if (this.innerX !== pos.x || this.y !== pos.y) {
        this.$emit('move', this.i, pos.x, pos.y);
      }

      this.calculateStyle();
    },

    dragoverHandler(event) {
      event.preventDefault();

      if (!this.dragging || this.resizing) {
        return;
      }

      const position = getControlPosition(event, this.layoutElement);
      const { x, y } = position;
      const coreEvent = createCoreData(this.lastX, this.lastY, x, y);

      this.dragging = {
        top: this.dragging.top + coreEvent.deltaY,
        left: this.dragging.left + coreEvent.deltaX,
      };

      const pos = this.calculateXY(this.dragging.top, this.dragging.left);

      this.lastX = x;
      this.lastY = y;

      if (this.innerX !== pos.x || this.y !== pos.y) {
        this.$emit('move', this.i, pos.x, pos.y);
      }

      this.calculateStyle();
    },

    dragendHandler(event) {
      if (!this.dragging || this.resizing) {
        return;
      }

      this.layoutElement.removeEventListener('dragover', this.throttledDragoverHandler);
      this.$el.removeEventListener('dragend', this.dragendHandler);

      const position = getControlPosition(event, this.layoutElement);
      const { x, y } = position;
      const coreEvent = createCoreData(this.lastX, this.lastY, x, y);
      const newPosition = {
        top: this.dragging.top + coreEvent.deltaY,
        left: this.dragging.left + coreEvent.deltaX,
      };

      this.dragging = null;

      const pos = this.calculateXY(newPosition.top, newPosition.left);

      this.$emit('moved', this.i, pos.x, pos.y);
      this.calculateStyle();
    },
  },
};
</script>

<style lang="scss" scoped>
.c-grid-item {
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

  &--resizing {
    opacity: .7;
    user-select: none;
  }

  &:not(&--disabled) {
    position: absolute;
    left: 0;
    right: auto;
    z-index: 2;
    pointer-events: none;
    overflow: hidden;
    height: 1000px;

    .c-grid-item__resize-handler {
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

    & ::v-deep * {
      pointer-events: none;
      overflow: hidden !important;
    }

    & ::v-deep .c-grid-item__resize-handler {
      &, * {
        pointer-events: auto;
      }
    }
  }
}
</style>
