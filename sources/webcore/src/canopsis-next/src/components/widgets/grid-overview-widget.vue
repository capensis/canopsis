<template lang="pug">
  div.gird-overview
    div.grid-item(
      v-for="widget in tab.widgets",
      :key="widget._id",
      :style="getWidgetFlexStyle(widget)"
    )
      widget-wrapper(
        :widget="widget",
        :tab="tab",
        :updateTabMethod="updateTabMethod"
      )
</template>

<script>
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  components: {
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
  methods: {
    getWidgetFlexStyle({ gridParameters }) {
      const parameters = this.getGridByWindowSize(gridParameters);

      return {
        gridColumnStart: parameters.x + 1,
        gridColumnEnd: parameters.x + 1 + parameters.w,
        gridRowStart: parameters.y + 1,
        gridRowEnd: parameters.fixedHeight
          ? parameters.y + parameters.h + 1
          : parameters.y + 2,
      };
    },

    getGridByWindowSize(gridParameters = {}) {
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
    display: grid;
    grid-gap: 10px;
    grid-template-columns: repeat(12, [col-start] 1fr);
  }
  .grid-item {
    overflow: auto;
  }
</style>
