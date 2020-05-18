<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    grid-overview-widget(
      v-if="!isEditingMode",
      :tab="tab",
      :updateTabMethod="updateTabMethod"
    )
</template>

<script>
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';
import GridOverviewWidget from '@/components/widgets/grid-overview-widget.vue';

import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  components: {
    GridOverviewWidget,
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
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    updateTabMethod: {
      type: Function,
      default: () => () => {},
    },
  },
};
</script>

<style lang="scss" scoped>
  .full-screen {
    .hide-on-full-screen {
      display: none;
    }
  }
</style>
