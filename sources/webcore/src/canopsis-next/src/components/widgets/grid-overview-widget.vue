<template lang="pug">
  div.grid-overview
    grid-overview-item(
      v-for="widget in tab.widgets",
      :widget="widget",
      :key="widget._id"
    )
      slot(:widget="widget")
</template>

<script>
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
};
</script>

<style lang="scss" scoped>
  .grid-overview {
    padding: 20px;
    display: grid;
    column-gap: 20px;
    grid-template-columns: repeat(12, 1fr);
  }
</style>
