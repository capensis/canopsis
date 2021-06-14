<template lang="pug">
  v-card.grid-item(:style="overviewItemStyles")
    slot
</template>

<script>
import { MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP, WIDGET_GRID_ROW_HEIGHT } from '@/constants';

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

      return this.widget.grid_parameters[key];
    },

    overviewItemStyles() {
      return {
        margin: `${WIDGET_GRID_ROW_HEIGHT / 2}px 0`,
        gridColumnStart: this.gridParameters.x + 1,
        gridColumnEnd: this.gridParameters.x + 1 + this.gridParameters.w,
        gridRowStart: this.gridParameters.y + 1,
        gridRowEnd: this.gridParameters.y + this.gridParameters.h + 1,
        height: this.gridParameters.autoHeight ? 'auto' : `${WIDGET_GRID_ROW_HEIGHT * this.gridParameters.h}px`,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .grid-item {
    overflow: auto;
  }
</style>
