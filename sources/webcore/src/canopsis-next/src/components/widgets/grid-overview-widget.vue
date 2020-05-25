<template lang="pug">
  div.gird-overview(:style="gridWrapperStyle")
    grid-overview-item(
      v-for="widget in tab.widgets",
      :widget="widget",
      :key="widget._id"
    )
      slot(:widget="widget")
</template>

<script>
import GridOverviewItem from '@/components/widgets/partials/grid-overview-item.vue';

import sideBarMixin from '@/mixins/side-bar/side-bar';
import { WIDGET_GRID_ROW_HEIGHT } from '@/constants';

export default {
  components: {
    GridOverviewItem,
  },
  mixins: [
    sideBarMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
  computed: {
    gridWrapperStyle() {
      return {
        // TODO fix 1000
        gridTemplateRows: `repeat(${1000}, ${WIDGET_GRID_ROW_HEIGHT}px)`,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .gird-overview {
    padding: 10px;
    display: grid;
    grid-template-columns: repeat(12, [col-start] 1fr);
  }
</style>
