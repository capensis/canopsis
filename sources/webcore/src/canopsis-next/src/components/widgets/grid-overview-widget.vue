<template lang="pug">
  div.gird-overview(:style="gridWrapperStyle")
    grid-overview-item(
      v-for="widget in widgets",
      :widget="widget",
      :key="widget._id"
    )
      slot(:widget="widget")
</template>

<script>
import GridOverviewItem from '@/components/widgets/partials/grid-overview-item.vue';

import sideBarMixin from '@/mixins/side-bar/side-bar';

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
    widgets() {
      return this.tab.widgets;
    },
  },
  methods: {
    gridParameters(widget) {
      const { gridParameters } = widget;

      return {
        xl: gridParameters.desktop,
        l: gridParameters.desktop,
        t: gridParameters.tablet,
        m: gridParameters.mobile,
      }[this.$mq];
    },
  },
};
</script>

<style lang="scss" scoped>
  .gird-overview {
    padding: 20px;
    display: grid;
    column-gap: 20px;
    grid-template-columns: repeat(12, [col-start] 1fr);
  }
</style>
