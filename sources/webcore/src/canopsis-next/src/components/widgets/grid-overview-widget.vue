<template lang="pug">
  div.grid-overview(:style="overviewLayoutStyles")
    grid-overview-item(
      v-for="widget in tab.widgets",
      :widget="widget",
      :key="widget._id"
    )
      slot(:widget="widget")
</template>

<script>
import { WIDGET_GRID_ROW_HEIGHT, WIDGET_GRID_COLUMNS_COUNT } from '@/constants';

import GridOverviewItem from '@/components/widgets/partials/grid-overview-item.vue';

export default {
  components: {
    GridOverviewItem,
  },
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
  computed: {
    overviewLayoutStyles() {
      const rowHeightInPixels = `${WIDGET_GRID_ROW_HEIGHT}px`;

      return {
        padding: rowHeightInPixels,
        columnGap: rowHeightInPixels,
        gridTemplateColumns: `repeat(${WIDGET_GRID_COLUMNS_COUNT}, 1fr)`,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .grid-overview {
    display: grid;
  }
</style>
