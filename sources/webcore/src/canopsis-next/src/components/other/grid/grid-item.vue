<script>
import { GridItem } from 'vue-grid-layout';

export default {
  extends: GridItem,
  methods: {
    autoSizeHeight() {
      // ok here we want to calculate if a resize is needed
      this.previousW = this.innerW;
      this.previousH = this.innerH;

      const newSize = this.$slots.default[0].elm.getBoundingClientRect();
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

    calcPosition(x, y, w, h) {
      const [marginX, marginY] = this.margin;
      const colWidth = this.calcColWidth();

      const result = {
        top: Math.round((this.rowHeight * y) + marginY),
        width: w === Infinity ? w : Math.round((colWidth * w) + (Math.max(0, w - 1) * marginX)),
        height: h === Infinity ? h : Math.round(this.rowHeight * h),
      };

      const horizontal = Math.round((colWidth * x) + ((x + 1) * marginX));

      if (this.renderRtl) {
        result.right = horizontal;
      } else {
        result.left = horizontal;
      }

      return result;
    },
  },
};
</script>
