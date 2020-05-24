<script>
import { GridItem } from 'vue-grid-layout';

export default {
  extends: GridItem,
  props: {
    fixedHeight: {
      type: Boolean,
      default: false,
    },
  },
  mounted() {
    if (!this.fixedHeight) {
      this.autoSizeHeight();
    }
  },
  methods: {
    autoSizeHeight() {
      // ok here we want to calculate if a resize is needed
      this.previousW = this.innerW;
      this.previousH = this.innerH;
      const defaultSlots = this.$slots.default;

      if (!defaultSlots) {
        return;
      }

      const newSize = defaultSlots[0].elm.getBoundingClientRect();
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
        this.$emit('resized', this.i, pos.h, this.previousW, newSize.height, newSize.width, true);
        this.eventBus.$emit('resizeEvent', 'resizeend', this.i, this.innerX, this.innerY, pos.h, this.innerW);
      }
    },

    countAboveCard(currentCardX, currentCardY, currentCardWidth) {
      const { layout } = this.$parent;

      const { count } = layout
        .filter(({ y }) => y < currentCardY)
        .sort((a, b) => b.y - a.y)
        .reduce((acc, { y, x, w }) => {
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

    calcWH(height, width) {
      const [marginX] = this.margin;
      const colWidth = this.calcColWidth();

      let w = Math.round((width + marginX) / (colWidth + marginX));
      let h = Math.round((height) / this.rowHeight);

      w = Math.max(Math.min(w, this.cols - this.innerX), 0);
      h = Math.max(Math.min(h, this.maxRows - this.innerY), 0);

      return { w, h };
    },

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
