<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    grid-overview-widget(
      v-show="!isEditingMode",
      :tab="tab"
    )
      template(#default="{ widget }")
        widget-wrapper(:widget="widget", :tab="tab")
    grid-edit-widgets(
      v-if="isEditingMode",
      :tab="tab",
      :update-tab-method="updateTabMethod",
      @update:widgets-fields="$emit('update:widgets-fields', $event)"
    )
      template(#default="{ widget }")
        widget-wrapper(:widget="widget", :tab="tab", is-editing-mode)
</template>

<script>
import GridOverviewWidget from '@/components/widgets/grid-overview-widget.vue';
import GridEditWidgets from '@/components/widgets/grid-edit-widgets.vue';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import queryMixin from '@/mixins/query';

export default {
  components: {
    WidgetWrapper,
    GridOverviewWidget,
    GridEditWidgets,
  },
  mixins: [
    queryMixin,
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
  destroyed() {
    this.removeWidgetsQueries();
  },
  methods: {
    /**
     * Remove queries which was created for all widgets
     */
    removeWidgetsQueries() {
      this.tab.widgets.forEach(({ _id: id }) => this.removeQuery({
        id,
      }));
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
