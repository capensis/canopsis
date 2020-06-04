<template lang="pug">
  v-card.grid-item(:style="overviewItemStyles")
    slot
</template>

<script>
import { MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP } from '@/constants';

export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    gridParameters() {
      const key = MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP[this.$mq];

      return this.widget.gridParameters[key];
    },

    overviewItemStyles() {
      return {
        gridColumnStart: this.gridParameters.x + 1,
        gridColumnEnd: this.gridParameters.x + 1 + this.gridParameters.w,
        gridRowStart: this.gridParameters.y + 1,
        gridRowEnd: this.gridParameters.y + this.gridParameters.h + 1,
        height: this.gridParameters.fixedHeight ? `${20 * this.gridParameters.h}px` : 'auto',
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .grid-item {
    overflow: auto;
    margin: 10px 0;
  }
</style>
