<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    component(
      :is="layoutComponent",
      :tab="tab",
      :isEditingMode="isEditingMode",
      :updateTabMethod="updateTabMethod"
    )
      widget-wrapper(
        slot-scope="props",
        :widget="props.widget",
        :tab="tab",
        :updateTabMethod="updateTabMethod"
      )
</template>

<script>
import GridOverviewWidget from '@/components/widgets/grid-overview-widget.vue';
import GridEditWidgets from '@/components/widgets/grid-edit-widgets.vue';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  components: {
    WidgetWrapper,
    GridOverviewWidget,
    GridEditWidgets,
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
  computed: {
    layoutComponent() {
      return this.isEditingMode ? 'grid-edit-widgets' : 'grid-overview-widget';
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
