<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    widget-wrapper(
      v-for="widget in tab.widgets",
      :key="widget._id",
      :widget="widget",
      :tab="tab",
      :isEditingMode="isEditingMode",
      :row="row",
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
