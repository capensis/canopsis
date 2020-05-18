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
  computed: {
    wrapperStyles() {
      return {
        gridTemplateRows: this.tab.widgets.reduce((style, { gridParameters }) => {
          const { h, fixedHeight } = gridParameters.desktop;

          return fixedHeight ? `${style} repeat(${h}, 20px)` : `${style} auto`;
        }, ''),
      };
    },
  },
  methods: {
    getWidgetFlexStyle(widget) {
      const parameters = this.getGridByWindowSize(widget.gridParameters);

      return {
        gridColumnStart: parameters.x + 1,
        gridColumnEnd: parameters.x + 1 + parameters.w,
        gridRowStart: parameters.y + 1,
        gridRowEnd: parameters.fixedHeight
          ? parameters.y + parameters.h
          : parameters.y + 2,
      };
    },

    getGridByWindowSize(gridParameters = {}) {
      return {
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
    grid-template-columns: repeat(12, calc(80% / 12));
  }
  .grid-item {
    overflow: auto;
  }
</style>
