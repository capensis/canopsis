<template lang="pug">
  div.gird-overview
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
};
</script>

<style lang="scss" scoped>
  .gird-overview {
    padding: 10px;
    display: grid;
    grid-template-columns: repeat(12, [col-start] 1fr);
  }
</style>
