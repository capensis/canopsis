<template lang="pug">
  div.c-grid-layout(ref="layout", :style="style")
    slot(:bind="bind", :on="on")
</template>

<script>
import { uniq } from 'lodash';

import { formArrayMixin } from '@/mixins/form';

const calculateLayoutBottom = (layout = []) => layout.reduce((acc, item) => Math.max(acc, item.y + item.h), 0);
const calculateLayoutRowCount = (layout = []) => uniq(layout.map(({ y }) => y)).length;

export default {
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
      dragging: false,
      resizing: false,
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
      };
    },
    on() {
      return {
        resized: this.resizedItemHandler,
      };
    },
  },
  mounted() {
    this.addAllListeners();
    this.resizeObserverHandler();
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

    resizedItemHandler(id, x, y, h, w) {
      const index = this.layout.findIndex(item => item.i === id);

      this.updateItemInArray(index, { ...this.layout[index], h, w });

      this.$nextTick(() => {
        this.updateHeight();
        this.resizing = false;
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
  background-color: rgba(60, 60, 60, 0.05);
}
</style>
