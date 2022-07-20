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
      @update:widgets-grid="updateWidgetsGrid"
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
      widgetsGrid: {},
    };
  },
  created() {
    this.registerEditingOffHandler(this.updatePositions);
  },
  beforeDestroy() {
    this.unregisterEditingOffHandler(this.updatePositions);
    this.removeWidgetsQueries();
  },
  methods: {
    updateWidgetsGrid(widgetsGrid) {
      this.widgetsGrid = widgetsGrid;
    },

    async updatePositions() {
      try {
        const data = Object.entries(this.widgetsGrid)
          .map(([id, gridParameters]) => ({
            _id: id,
            grid_parameters: gridParameters,
          }));

        if (!data.length) {
          return;
        }

        await this.updateWidgetGridPositions({ data });
        await this.fetchActiveView();
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
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
