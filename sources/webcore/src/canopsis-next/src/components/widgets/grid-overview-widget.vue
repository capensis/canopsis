<template lang="pug">
  div.gird-overview
    grid-overview-item(
      v-for="widget in tab.widgets",
      :widget="widget",
      :key="widget._id"
    )
      widget-wrapper(
        :widget="widget",
        :tab="tab",
        :updateTabMethod="updateTabMethod"
      )
</template>

<script>
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';
import GridOverviewItem from '@/components/widgets/partials/grid-overview-item.vue';

import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  components: {
    GridOverviewItem,
    WidgetWrapper,
  },
  mixins: [
    sideBarMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
    updateTabMethod: {
      type: Function,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
  .gird-overview {
    padding: 10px;
    display: grid;
    grid-gap: 10px;
    grid-template-columns: repeat(12, [col-start] 1fr);
  }
</style>
