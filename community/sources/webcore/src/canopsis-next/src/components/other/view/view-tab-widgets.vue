<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    grid-overview-widget(
      v-show="!editing",
      :tab="tab"
    )
      template(#default="{ widget }")
        widget-wrapper(:widget="widget", :tab="tab")
    grid-edit-widgets(
      v-if="editing",
      :tab="tab",
      @update:widgets-fields="updateWidgetsFieldsForUpdateById"
    )
      template(#default="{ widget }")
        widget-wrapper(:widget="widget", :tab="tab", editing)
</template>

<script>
import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

import GridOverviewWidget from '@/components/widgets/grid-overview-widget.vue';
import GridEditWidgets from '@/components/widgets/grid-edit-widgets.vue';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

export default {
  components: {
    WidgetWrapper,
    GridOverviewWidget,
    GridEditWidgets,
  },
  mixins: [
    queryMixin,
    activeViewMixin,
    entitiesWidgetMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      widgetsFieldsForUpdateById: {},
    };
  },
  created() {
    this.registerEditingOffHandler(this.updatePositions);
  },
  beforeDestroy() {
    this.removeWidgetsQueries();
    this.unregisterEditingOffHandler(this.updatePositions);
  },
  methods: {
    updateWidgetsFieldsForUpdateById(widgetsFieldsForUpdateById) {
      this.widgetsFieldsForUpdateById = {
        ...this.widgetsFieldsForUpdateById,
        ...widgetsFieldsForUpdateById,
      };
    },

    updatePositions() {
      return new Promise(resolve => setTimeout(resolve, 5000));
    },

    /**
     * Remove queries which was created for all widgets
     */
    removeWidgetsQueries() {
      this.tab.widgets.forEach(({ _id: id }) => this.removeQuery({ id }));
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
